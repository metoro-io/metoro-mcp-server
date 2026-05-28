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
	AlertName         string                  `json:"alert_name" jsonschema:"required,description=The name of the alert to create"`
	AlertDescription  string                  `json:"alert_description" jsonschema:"required,description=The description of the alert to create"`
	Timeseries        []model.MetricSpecifier `json:"timeseries" jsonschema:"required,description=Array of timeseries data to get. Each item in this array corresponds to a single timeseries. You can then use the formulas to combine these timeseries. If you only want to see the combination of timeseries via defining formulas and if you dont want to see the individual timeseries data when setting formulas you can set shouldNotReturn to true. For each timeseries make sure to set the type."`
	Formula           model.Formula           `json:"formula" jsonschema:"description=Optional formula to combine timeseries. Formula should only consist of formulaIdentifier of the timeseries in the timeseries array. e.g. a + b + c if a b c appears in the formulaIdentifier of the timeseries array. You can ONLY do the following operations: Arithmetic operations:+ (for add) - (for substract) * (for multiply) / (for division) % (for modulus) ^ or ** (for exponent). Comparison: == != < > <= >= . Logical:! (for not) && (for AND) || (for OR). Conditional operations: ?: (ternary) e.g. (a || b) ? 1 : 0. Do not guess the operations. Just use these available ones!"`
	Condition         string                  `json:"condition" jsonschema:"required,enum=GreaterThan,enum=LessThan,enum=GreaterThanOrEqual,enum=LessThanOrEqual,description=the arithmetic comparison to use to evaluate whether an alert is firing or not. This is used to determine whether the alert should be triggered based on the threshold value."`
	Threshold         float64                 `json:"threshold" jsonschema:"required,description=The threshold value for the alert. This is the value that will be used together with the the arithmetic condition to see whether the alert should be triggered or not. For example if you set the condition to GreaterThan and the threshold to 100 then the alert will fire if the value of the timeseries is greater than 100."`
	DatapointsToAlarm int64                   `json:"datapoints_to_alarm" jsonschema:"required,description=The number of datapoints that need to breach the threshold for the alert to be triggered"`
	EvaluationWindow  int64                   `json:"evaluation_window" jsonschema:"required,description=The evaluation window in number of datapoints. This is the number of datapoints that will be considered for evaluating the alert condition. For example if you set this to then the last 5 datapoints will be considered for evaluating the alert condition. This is useful for smoothing out spikes in the data and preventing false positives."`
}

func CreateAlertHandler(ctx context.Context, arguments CreateAlertHandlerArgs) (*mcpgolang.ToolResponse, error) {
	alert, err := createAlertFromTimeseries(ctx, arguments.AlertName, arguments.AlertDescription, arguments.Timeseries, arguments.Formula, arguments.Condition, arguments.Threshold, arguments.DatapointsToAlarm, arguments.EvaluationWindow)
	if err != nil {
		return nil, fmt.Errorf("error creating alert properties: %v", err)
	}

	newAlertRequest := model.CreateUpdateAlertRequest{
		Alert: alert,
	}

	resp, err := setAlertMetoroCall(ctx, newAlertRequest)
	if err != nil {
		return nil, fmt.Errorf("error setting dashboard: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

// TODO: Implement the conversion logic.
func createAlertFromTimeseries(ctx context.Context, alertName, alertDescription string, timeseries []model.MetricSpecifier, formula model.Formula, condition string, threshold float64, datapointsToAlarm int64, evaluationWindow int64) (model.Alert, error) {
	// Create dummy time range for the last 10 minutes to validate the timeseries
	endTime := time.Now().Unix()
	startTime := endTime - 600 // 10 minutes ago

	// Convert MetricSpecifier to SingleTimeseriesRequest for validation
	singleTimeseriesRequests := convertMetricSpecifierToSingleTimeseries(timeseries)

	err := checkTimeseries(ctx, singleTimeseriesRequests, startTime, endTime)
	if err != nil {
		return model.Alert{}, err
	}
	metoroQlQueries, err := convertMetricSpecifierToMetoroQL(ctx, timeseries, []model.Formula{formula})
	if err != nil {
		return model.Alert{}, fmt.Errorf("error converting metric specifiers to MetoroQL: %v", err)
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

	// Determine bucket size from the timeseries
	bucketSize := int64(60) // default to 60 seconds
	if len(timeseries) > 0 && timeseries[0].BucketSize > 0 {
		bucketSize = timeseries[0].BucketSize
	}

	// Use the first MetoroQL query (usually the combined formula result)
	query := ""
	if len(metoroQlQueries) > 0 {
		for _, q := range metoroQlQueries {
			if q != "" {
				query = q
				break
			}
		}
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
					Query:      query,
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

	return alert, nil
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
			ShouldNotReturn:   spec.ShouldNotReturn,
			FormulaIdentifier: "",
		}
	}
	return result
}
