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
	TimeConfig utils.TimeConfig                `json:"time_config" jsonschema:"required,description=The time period to get the timeseries data for. e.g. if you want to get the timeseries data for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Timeseries []model.SingleTimeseriesRequest `json:"timeseries" jsonschema:"required,description=Array of timeseries data to get. Each item in this array corresponds to a single timeseries. You can then use the formulas to combine these timeseries. If you only want to see the combination of timeseries via defining formulas and if you dont want to see the individual timeseries data when setting formulas you can set shouldNotReturn to true"`
	Formulas   []model.Formula                 `json:"formulas" jsonschema:"description=Optional formulas to combine timeseries. Formula should only consist of formulaIdentifier of the timeseries in the timeseries array. e.g. a + b + c if a b c appears in the formulaIdentifier of the timeseries array. You can ONLY do the following operations: Arithmetic operations:+ (for add) - (for substract) * (for multiply) / (for division) % (for modulus) ^ or ** (for exponent). Comparison: == != < > <= >= . Logical:! (for not) && (for AND) || (for OR). Conditional operations: ?: (ternary) e.g. (a || b) ? 1 : 0. Do not guess the operations. Just use these available ones!"`
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

func convertTimeseriesToAPITimeseries(timeseries []model.SingleTimeseriesRequest, startTime int64, endTime int64) []model.SingleMetricRequest {
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
		case model.KubernetesResource:
			apiRequest.KubernetesResource = &model.GetKubernetesResourceRequest{
				StartTime:      startTime,
				EndTime:        endTime,
				Filters:        ts.Filters,
				ExcludeFilters: ts.ExcludeFilters,
				Splits:         ts.Splits,
				BucketSize:     ts.BucketSize,
				Functions:      ts.Functions,
				JsonPath:       ts.JsonPath,
				Aggregation:    ts.Aggregation,
			}
		}
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

func checkTimeseries(ctx context.Context, timeseries []model.SingleTimeseriesRequest, startTime, endTime int64) error {
	for _, ts := range timeseries {
		switch ts.Type {
		case model.Metric:
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
		case model.Trace:
			err := CheckAttributes(ctx, ts.Type, ts.Filters, ts.ExcludeFilters, ts.Splits, nil)
			if err != nil {
				return err
			}
		case model.Logs:
			err := CheckAttributes(ctx, ts.Type, ts.Filters, ts.ExcludeFilters, ts.Splits, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
