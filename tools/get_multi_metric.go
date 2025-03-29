package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"

	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetMultiMetricHandlerArgs struct {
	TimeConfig utils.TimeConfig      `json:"time_config" jsonschema:"required,description=The time period to get the timeseries data for. e.g. if you want to get the timeseries data for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Timeseries []SingleMetricRequest `json:"timeseries" jsonschema:"required,description=Array of timeseries data descriptions to get."`
	Formulas   []model.Formula       `json:"formulas" jsonschema:"description=Optional formulas to combine timeseries. Formula should only consist of formulaIdentifier of the timeseries in the timeseries array. e.g. a + b + c if a b c appears in the formulaIdentifier of the timeseries array. You can do * / + - && || ?: operations"`
}

// TODO: Add kubernetes resource request type attributes.
type SingleMetricRequest struct {
	Type              model.MetricType    `json:"type" jsonschema:"required,enum=metric,enum=trace,enum=logs,enum=kubernetes_resource,description=Type of timeseries data to retrieve"`
	MetricName        string              `json:"metricName" jsonschema:"description=THIS IS ONLY REQUIRED IF THE type is 'metric'.The name of the metric to use for getting the timeseries data for type 'metric'"`
	Aggregation       model.Aggregation   `json:"aggregation" jsonschema:"required,enum=sum,enum=count,enum=min,enum=max,enum=avg,enum=p50,enum=p90,enum=p95,enum=p99,description=The aggregation to apply to the timeseries at the datapoint bucket size level. The aggregation will be applied to every datapoint bucket. For example if the bucket size is 1 minute and the aggregation is sum then the sum of all datapoints in a minute will be returned"`
	Filters           map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the timeseries. Only the timeseries that match these filters will be returned. Get the possible filter keys and values from the get_attribute_keys and get_attribute_values tools. e.g. {service_name: [/k8s/namespaceX/serviceX]} should return timseries for serviceX in namespaceX. This is just and example. Do not guess the attribute keys and values."`
	ExcludeFilters    map[string][]string `json:"excludeFilters" jsonschema:"description=Filters to exclude the timeseries data. Timeseries matching the exclude filters will not be returned. Get the possible exclude filter keys and values from the get_attribute_keys and get_attribute_values tools. e.g. {service_name: [/k8s/namespaceX/serviceX]} should exclude timeseries from serviceX in namespaceX. This is just and example. Do not guess the attribute keys and values"`
	Splits            []string            `json:"splits" jsonschema:"description=Splits will allow you to group/split timeseries data by an attribute. This is useful if you would like to see the breakdown of a particular timeseries by an attribute. For example if you want to see the breakdown of the metric by attributeX you would set the splits as ['attributeX']. Get the attribute keys from the get_attribute_keys tool. Do not guess"`
	Regexes           []string            `json:"regexes" jsonschema:"description=This should only be set if the type is 'logs'. Regexes are evaluated against the log message/body. Regexes to apply to the timeseries data. Only the timeseries data that match these regexes will be returned. Regexes are ORed together. For example if you want to get timeseries data with message that contains the word 'error' or 'warning' you would set the regexes as ['error' 'warning']"`
	ExcludeRegexes    []string            `json:"excludeRegexes" jsonschema:"description=This should only be set if the type is 'logs'. Regexes are evaluated against the log message/body. Regexes to exclude the timeseries data. Timeseries data that match these regexes will not be returned. Exclude regexes are AND together. For example if you want to get timeseries data with messages that do not contain the word 'error' or 'warning' you would set the exclude regexes as ['error' 'warning']"`
	BucketSize        int64               `json:"bucketSize" jsonschema:"description=The size of each datapoint bucket in seconds if not provided metoro will select the best bucket size for the given duration for performance and clarity"`
	ShouldNotReturn   bool                `json:"shouldNotReturn" jsonschema:"description=If true result won't be returned (useful for formulas)"`
	FormulaIdentifier string              `json:"formulaIdentifier" jsonschema:"description=Identifier to reference this metric in formulas"`
}

func GetMultiMetricHandler(ctx context.Context, arguments GetMultiMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetMultiMetricRequest{
		StartTime: startTime,
		EndTime:   endTime,
		Metrics:   convertTimeseriesToAPITimeseries(arguments.Timeseries),
		Formulas:  arguments.Formulas,
	}

	if len(arguments.Timeseries) == 0 {
		return nil, fmt.Errorf("no timeseries data provided")
	}

	body, err := getMultiMetricMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting metric: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getMultiMetricMetoroCall(ctx context.Context, request model.GetMultiMetricRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling metric request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "metrics", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}

func convertTimeseriesToAPITimeseries(timeseries []SingleMetricRequest) []model.SingleMetricRequest {
	result := make([]model.SingleMetricRequest, len(timeseries))

	for i, ts := range timeseries {
		apiRequest := model.SingleMetricRequest{
			Type:              string(ts.Type),
			ShouldNotReturn:   ts.ShouldNotReturn,
			FormulaIdentifier: ts.FormulaIdentifier,
		}

		switch ts.Type {
		case model.Metric:
			apiRequest.Metric = &model.GetMetricRequest{
				MetricName:     ts.MetricName,
				Filters:        ts.Filters,
				ExcludeFilters: ts.ExcludeFilters,
				Splits:         ts.Splits,
				Aggregation:    ts.Aggregation,
				//Functions:      ts.Metric.Functions,
				//LimitResults:   ts.Metric.LimitResults,
				BucketSize: ts.BucketSize,
			}

		case model.Trace:
			apiRequest.Trace = &model.GetTraceMetricRequest{
				Filters:        ts.Filters,
				ExcludeFilters: ts.ExcludeFilters,
				Splits:         ts.Splits,
				Aggregate:      ts.Aggregation,
				BucketSize:     ts.BucketSize,
				//ServiceNames:   ts.ServiceNames,
				//Regexes:        ts.Regexes,
				//ExcludeRegexes: ts.ExcludeRegexes,
				//Environments:   ts.Environments,
				//LimitResults:   ts.LimitResults,
			}

		case model.Logs:
			apiRequest.Logs = &model.GetLogMetricRequest{
				GetLogsRequest: model.GetLogsRequest{
					Filters:        ts.Filters,
					ExcludeFilters: ts.ExcludeFilters,
					Regexes:        ts.Regexes,
					ExcludeRegexes: ts.ExcludeRegexes,
					//Environments:   ts.Environments,
				},
				Splits:     ts.Splits,
				BucketSize: ts.BucketSize,
				//Functions:  ts.Functions,
			}
		}

		// TODO: Add the kubernetes resources here.

		result[i] = apiRequest
	}

	return result
}
