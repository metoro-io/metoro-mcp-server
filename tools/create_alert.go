package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type CreateAlertHandlerArgs struct {
	AlertName         string  `json:"alert_name" jsonschema:"required,description=The name of the alert to create"`
	AlertDescription  string  `json:"alert_description" jsonschema:"required,description=The description of the alert to create"`
	MetoroQL          string  `json:"metoroql" jsonschema:"required,description=The MetoroQL query to evaluate for this alert. Use get_metric_names, get_attribute_keys, get_attribute_values, and get_timeseries_data to discover and validate query inputs before creating an alert."`
	BucketSize        int64   `json:"bucket_size" jsonschema:"description=The size of each datapoint bucket in seconds. Defaults to 60 seconds when omitted."`
	Condition         string  `json:"condition" jsonschema:"required,enum=GreaterThan,enum=LessThan,enum=GreaterThanOrEqual,enum=LessThanOrEqual,description=the arithmetic comparison to use to evaluate whether an alert is firing or not. This is used to determine whether the alert should be triggered based on the threshold value."`
	Threshold         float64 `json:"threshold" jsonschema:"required,description=The threshold value for the alert. This is the value that will be used together with the the arithmetic condition to see whether the alert should be triggered or not. For example if you set the condition to GreaterThan and the threshold to 100 then the alert will fire if the value of the timeseries is greater than 100."`
	DatapointsToAlarm int64   `json:"datapoints_to_alarm" jsonschema:"required,description=The number of datapoints that need to breach the threshold for the alert to be triggered"`
	EvaluationWindow  int64   `json:"evaluation_window" jsonschema:"required,description=The evaluation window in number of datapoints. This is the number of datapoints that will be considered for evaluating the alert condition. For example if you set this to then the last 5 datapoints will be considered for evaluating the alert condition. This is useful for smoothing out spikes in the data and preventing false positives."`
	InvestigateOnFire bool    `json:"investigate_on_fire" jsonschema:"description=Whether Metoro should automatically start an AI investigation when this alert fires. Defaults to false."`
}

func CreateAlertHandler(ctx context.Context, arguments CreateAlertHandlerArgs) (*mcpgolang.ToolResponse, error) {
	alert, err := createAlertFromMetoroQL(ctx, arguments.AlertName, arguments.AlertDescription, arguments.MetoroQL, arguments.BucketSize, arguments.Condition, arguments.Threshold, arguments.DatapointsToAlarm, arguments.EvaluationWindow, arguments.InvestigateOnFire)
	if err != nil {
		return nil, fmt.Errorf("error creating alert properties: %v", err)
	}

	newAlertRequest := model.CreateUpdateAlertRequest{
		Alert: alert,
	}

	resp, err := setAlertMetoroCall(ctx, newAlertRequest)
	if err != nil {
		return nil, fmt.Errorf("error setting alert: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func createAlertFromMetoroQL(ctx context.Context, alertName, alertDescription, metoroQL string, bucketSize int64, condition string, threshold float64, datapointsToAlarm int64, evaluationWindow int64, investigateOnFire bool) (model.Alert, error) {
	if metoroQL == "" {
		return model.Alert{}, fmt.Errorf("metoroql is required")
	}
	if bucketSize == 0 {
		bucketSize = 60
	}
	if bucketSize < 0 {
		return model.Alert{}, fmt.Errorf("bucket_size must be positive")
	}

	if err := validateCreateAlertMetoroQL(ctx, metoroQL, bucketSize); err != nil {
		return model.Alert{}, err
	}

	// Convert condition string to OperatorType
	var operatorType model.OperatorType
	switch condition {
	case "GreaterThan":
		operatorType = model.GREATER_THAN
	case "LessThan":
		operatorType = model.LESS_THAN
	case "GreaterThanOrEqual":
		operatorType = model.GREATER_THAN_OR_EQUAL
	case "LessThanOrEqual":
		operatorType = model.LESS_THAN_OR_EQUAL
	default:
		return model.Alert{}, fmt.Errorf("invalid condition: %s", condition)
	}

	// Create the alert
	conditionType := model.STATIC
	timeseriesType := model.TIMESERIES
	alert := model.Alert{
		Metadata: model.MetadataObject{
			Name:        alertName,
			Description: &alertDescription,
			Id:          uuid.NewString(),
		},
		Type: &timeseriesType,
		Timeseries: model.TimeseriesConfig{
			Expression: model.ExpressionConfig{
				MetoroQLTimeseries: &model.MetoroQlTimeseries{
					Query:      metoroQL,
					BucketSize: bucketSize,
				},
			},
			EvaluationRules: []model.Condition{
				{
					Name: "Alert Condition",
					Type: &conditionType,
					Static: &model.StaticCondition{
						Operators: []model.OperatorConfig{
							{
								Operator:  operatorType,
								Threshold: threshold,
							},
						},
						PersistenceSettings: model.PersistenceSettings{
							DatapointsToAlarm:            datapointsToAlarm,
							DatapointsInEvaluationWindow: evaluationWindow,
						},
					},
				},
			},
		},
	}
	if investigateOnFire {
		alert.SetInvestigateOnFire(true)
	}

	return alert, nil
}

type metoroQLQueriesRequest struct {
	Queries []string `json:"queries"`
}

type metoroQLQueryResponse struct {
	MetricSpecifiers []model.MetricSpecifier `json:"metricSpecifiers"`
	Formulas         []model.Formula         `json:"formulas"`
}

func validateCreateAlertMetoroQL(ctx context.Context, metoroQL string, bucketSize int64) error {
	metricSpecifiers, err := convertMetoroQLToMetricSpecifiers(ctx, metoroQL)
	if err != nil {
		return fmt.Errorf("invalid MetoroQL query: %v", err)
	}

	endTime := time.Now().Unix()
	startTime := endTime - 3600

	singleTimeseriesRequests := convertMetricSpecifierToSingleTimeseries(metricSpecifiers)
	for i := range singleTimeseriesRequests {
		if singleTimeseriesRequests[i].BucketSize == 0 {
			singleTimeseriesRequests[i].BucketSize = bucketSize
		}
	}

	if err := checkTimeseries(ctx, singleTimeseriesRequests, startTime, endTime); err != nil {
		return fmt.Errorf("invalid MetoroQL query: %v", err)
	}
	return nil
}

func convertMetoroQLToMetricSpecifiers(ctx context.Context, metoroQL string) ([]model.MetricSpecifier, error) {
	req := metoroQLQueriesRequest{Queries: []string{metoroQL}}
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling MetoroQL conversion request: %v", err)
	}
	resp, err := utils.MakeMetoroAPIRequest("POST", "metoroql/convert/metoroqlToMetricSpecifier", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making MetoroQL conversion request: %v", err)
	}
	var metoroQLResp metoroQLQueryResponse
	if err := json.Unmarshal(resp, &metoroQLResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling MetoroQL conversion response: %v", err)
	}
	if len(metoroQLResp.MetricSpecifiers) == 0 {
		return nil, fmt.Errorf("no metric specifiers returned from MetoroQL conversion")
	}
	return metoroQLResp.MetricSpecifiers, nil
}

func convertMetricSpecifierToMetoroQL(ctx context.Context, metricSpecs []model.MetricSpecifier, formulas []model.Formula) ([]string, error) {
	req := model.MetricSpecifiersRequest{
		MetricSpecifiers: metricSpecs,
		Formulas:         formulas,
	}
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling MetricSpecifiersRequest: %v", err)
	}
	resp, err := utils.MakeMetoroAPIRequest("POST", "metoroql/convert/metricSpecifierToMetoroql", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making MetoroQL conversion request: %v", err)
	}
	var metoroQLQueriesResp model.MetricSpecifierToMetoroQLResponse
	if err := json.Unmarshal(resp, &metoroQLQueriesResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling MetoroQL conversion response: %v", err)
	}
	if len(metoroQLQueriesResp.Queries) == 0 {
		return nil, fmt.Errorf("no MetoroQL queries returned from conversion")
	}
	return metoroQLQueriesResp.Queries, nil
}

func setAlertMetoroCall(ctx context.Context, request model.CreateUpdateAlertRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling alert request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "alerts/update", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}

// convertMetricSpecifierToSingleTimeseries converts MetricSpecifier to SingleTimeseriesRequest
func convertMetricSpecifierToSingleTimeseries(metricSpecs []model.MetricSpecifier) []model.SingleTimeseriesRequest {
	result := make([]model.SingleTimeseriesRequest, len(metricSpecs))
	for i, spec := range metricSpecs {
		result[i] = model.SingleTimeseriesRequest{
			Type:              spec.MetricType,
			MetricName:        spec.MetricName,
			Aggregation:       spec.Aggregation,
			Filters:           model.MapToFilters(spec.Filters),
			ExcludeFilters:    model.MapToFilters(spec.ExcludeFilters),
			Splits:            spec.Splits,
			Regexes:           spec.Regexes,
			ExcludeRegexes:    spec.ExcludeRegexes,
			BucketSize:        spec.BucketSize,
			Functions:         spec.Functions,
			JsonPath:          spec.JsonPath,
			ShouldNotReturn:   spec.ShouldNotReturn,
			FormulaIdentifier: "",
		}
	}
	return result
}
