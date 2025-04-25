package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"

	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetMultiMetricHandlerArgs struct {
	TimeConfig utils.TimeConfig      `json:"time_config" jsonschema:"required,description=The time period to get the timeseries data for. e.g. if you want to get the timeseries data for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Timeseries []SingleMetricRequest `json:"timeseries" jsonschema:"required,description=Array of timeseries data to get. Each item in this array corresponds to a single timeseries. You can then use the formulas to combine these timeseries. If you only want to see the combination of timeseries via defining formulas and if you dont want to see the individual timeseries data when setting formulas you can set shouldNotReturn to true"`
	Formulas   []model.Formula       `json:"formulas" jsonschema:"description=Optional formulas to combine timeseries. Formula should only consist of formulaIdentifier of the timeseries in the timeseries array. e.g. a + b + c if a b c appears in the formulaIdentifier of the timeseries array. You can ONLY do the following operations: Arithmetic operations:+ (for add) - (for substract) * (for multiply) / (for division) % (for modulus) ^ or ** (for exponent). Comparison: == != < > <= >= . Logical:! (for not) && (for AND) || (for OR). Conditional operations: ?: (ternary) e.g. (a || b) ? 1 : 0. Do not guess the operations. Just use these available ones!"`
}

// TODO: Add kubernetes resource request type attributes.
type SingleMetricRequest struct {
	Type              model.MetricType       `json:"type" jsonschema:"required,enum=metric,enum=trace,enum=logs,enum=kubernetes_resource,description=Type of timeseries data to retrieve"`
	MetricName        string                 `json:"metricName" jsonschema:"description=THIS IS ONLY REQUIRED IF THE type is 'metric'.The name of the metric to use for getting the timeseries data for type 'metric'. If metric name ends with _total metoro already accounts for rate differences when returning the value so you don't need to calculate the rate yourself."`
	Aggregation       model.Aggregation      `json:"aggregation" jsonschema:"required,enum=sum,enum=count,enum=min,enum=max,enum=avg,enum=p50,enum=p90,enum=p95,enum=p99,description=The aggregation to apply to the timeseries at the datapoint bucket size level. The aggregation will be applied to every datapoint bucket. For example if the bucket size is 1 minute and the aggregation is sum then the sum of all datapoints in a minute will be returned. Do not guess the aggregations. Use the available ones. For traces you can use count p50 p90 p95 p99. for logs its always count. For metrics you can use sum min max avg"`
	Filters           map[string][]string    `json:"filters" jsonschema:"description=Filters to apply to the timeseries. Only the timeseries that match these filters will be returned. You MUST call get_attribute_keys and get_attribute_values tools to get the valid filter keys and values. e.g. {service_name: [/k8s/namespaceX/serviceX]} should return timeseries for serviceX in namespaceX. This is just and example. Do not guess the attribute keys and values."`
	ExcludeFilters    map[string][]string    `json:"excludeFilters" jsonschema:"description=Filters to exclude the timeseries data. Timeseries matching the exclude filters will not be returned. You MUST call get_attribute_keys and get_attribute_values tools to get the valid filter keys and values. e.g. {service_name: [/k8s/namespaceX/serviceX]} should exclude timeseries from serviceX in namespaceX. This is just and example. Do not guess the attribute keys and values"`
	Splits            []string               `json:"splits" jsonschema:"description=Array of attribute keys to split/group by the timeseries data by. Splits will allow you to group timeseries data by an attribute. This is useful if you would like to see the breakdown of a particular timeseries by an attribute. Get the attributes that you can pass into as Splits from the get_attribute_keys tool. DO NOT GUESS THE ATTRIBUTES."`
	Regexes           []string               `json:"regexes" jsonschema:"description=This should only be set if the type is 'logs'. Regexes are evaluated against the log message/body. Only the timeseries (logs) data that match these regexes will be returned. Regexes are ANDed together. For example if you want to get log count with message that contains the words 'fish' and 'chips' you would set the regexes as ['fish' 'chips']"`
	ExcludeRegexes    []string               `json:"excludeRegexes" jsonschema:"description=This should only be set if the type is 'logs'. Exclude regexes are evaluated against the log message/body. Log timeseries data that match these regexes will not be returned. Exclude regexes are ORed together. For example if you want to get timeseries data with messages that do not contain the word 'fish' or 'chips' you would set the exclude regexes as ['fish' 'chips']"`
	BucketSize        int64                  `json:"bucketSize" jsonschema:"description=The size of each datapoint bucket in seconds if not provided metoro will select the best bucket size for the given duration for performance and clarity"`
	Functions         []model.MetricFunction `json:"functions" jsonschema:"description=Array of functions to apply to the timeseries data in the order as it appears in the array. Functions will be applied to the timeseries data after the aggregation. For example if the aggregation is sum and the function is perSecond then the perSecond of the sum will be returned. Do not guess the functions. Use the available ones. For traces you can use rate. For logs you can use count. For metrics you can use rate sum min max avg. For kubernetes resources you can use rate sum min max avg"`
	ShouldNotReturn   bool                   `json:"shouldNotReturn" jsonschema:"description=If true result won't be returned (useful for formulas). Only set this to true if you only want to see the combination of timeseries via defining formulas and if you dont want to see the individual timeseries data.'"`
	FormulaIdentifier string                 `json:"formulaIdentifier" jsonschema:"description=Identifier to reference this metric in formulas. These should be unique for timeseries that you are requesting. For ease of use you can just use the alpahabet letters. e.g. a b c d e f g h i j k l m n o p q r s t u v w x y z"`
}

func GetMultiMetricHandler(ctx context.Context, arguments GetMultiMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	err = checkTimeseries(ctx, arguments.Timeseries, startTime, endTime)
	if err != nil {
		return nil, err
	}

	request := model.GetMultiMetricRequest{
		StartTime: startTime,
		EndTime:   endTime,
		Metrics:   convertTimeseriesToAPITimeseries(arguments.Timeseries, startTime, endTime),
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

func convertTimeseriesToAPITimeseries(timeseries []SingleMetricRequest, startTime int64, endTime int64) []model.SingleMetricRequest {
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
				StartTime:      startTime,
				EndTime:        endTime,
				MetricName:     ts.MetricName,
				Filters:        ts.Filters,
				ExcludeFilters: ts.ExcludeFilters,
				Splits:         ts.Splits,
				Aggregation:    ts.Aggregation,
				Functions:      ts.Functions,
				//Functions:      ts.Metric.Functions,
				//LimitResults:   ts.Metric.LimitResults,
				BucketSize: ts.BucketSize,
			}

		case model.Trace:
			apiRequest.Trace = &model.GetTraceMetricRequest{
				StartTime:      startTime,
				EndTime:        endTime,
				Filters:        ts.Filters,
				ExcludeFilters: ts.ExcludeFilters,
				Splits:         ts.Splits,
				Aggregate:      ts.Aggregation,
				BucketSize:     ts.BucketSize,
				Functions:      ts.Functions,
				//ServiceNames:   ts.ServiceNames,
				//Regexes:        ts.Regexes,
				//ExcludeRegexes: ts.ExcludeRegexes,
				//Environments:   ts.Environments,
				//LimitResults:   ts.LimitResults,
				//
			}

		case model.Logs:
			apiRequest.Logs = &model.GetLogMetricRequest{
				GetLogsRequest: model.GetLogsRequest{
					StartTime:      startTime,
					EndTime:        endTime,
					Filters:        ts.Filters,
					ExcludeFilters: ts.ExcludeFilters,
					Regexes:        ts.Regexes,
					ExcludeRegexes: ts.ExcludeRegexes,
					//Environments:   ts.Environments,
				},
				Functions:  ts.Functions,
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

func CheckAttributes(ctx context.Context, requestType model.MetricType, filters map[string][]string, excludeFilters map[string][]string, splits []string, metricRequest *model.GetMetricAttributesRequest) error {
	// Check whether the attributes given are valid.
	request := model.MultiMetricAttributeKeysRequest{
		Type:   string(requestType),
		Metric: metricRequest,
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	attributeResp, err := utils.MakeMetoroAPIRequest("POST", "metrics/attributes", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return fmt.Errorf("error making Metoro call: %v", err)
	}

	attributeKeys := model.GetAttributeKeysResponse{}
	err = json.Unmarshal(attributeResp, &attributeKeys)
	if err != nil {
		return fmt.Errorf("error unmarshaling response: %v", err)
	}

	attributesAsString := strings.Join(attributeKeys.Attributes, ", ")

	// Check whether the filters given are valid.
	for key, _ := range filters {
		if !slices.Contains(attributeKeys.Attributes, key) {
			return fmt.Errorf("invalid filter key: %s. Valid filter keys are: %s. Please try again with a valid key", key, attributesAsString)
		}
	}

	for key, _ := range excludeFilters {
		if !slices.Contains(attributeKeys.Attributes, key) {
			return fmt.Errorf("invalid exclude filter key: %s. Valid keys are: %s. Please try again with a valid key", key, attributesAsString)
		}
	}

	for _, split := range splits {
		if !slices.Contains(attributeKeys.Attributes, split) {
			return fmt.Errorf("invalid split key: %s. Valid keys are: %s. Please try again with a valid key", split, attributesAsString)
		}
	}
	return nil
}

func checkTimeseries(ctx context.Context, timeseries []SingleMetricRequest, startTime, endTime int64) error {
	for _, ts := range timeseries {
		if ts.Type != model.Metric {
			continue
		}
		err := CheckMetric(ctx, ts.MetricName)
		if err != nil {
			return err
		}
		err = CheckAttributes(ctx, ts.Type, ts.Filters, ts.ExcludeFilters, ts.Splits, &model.GetMetricAttributesRequest{
			StartTime:  startTime,
			EndTime:    endTime,
			MetricName: ts.MetricName,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
