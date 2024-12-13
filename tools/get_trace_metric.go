package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetTraceMetricHandlerArgs struct {
	TimeConfig     utils.TimeConfig       `json:"time_config" jsonschema:"required,description=The time period to get the trace timeseries data for. e.g. if you want to get the trace timeseries data for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	ServiceNames   []string               `json:"serviceNames" jsonschema:"description=Service names to return the trace timeseries data for"`
	Filters        map[string][]string    `json:"filters" jsonschema:"description=The filters to apply to the traces. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string    `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the traces. Traces matching the exclude traces will not be returned. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	Regexes        []string               `json:"regexes" jsonschema:"description=The regexes to apply to the trace's endpoints. Traces with endpoints matching regexes will be returned"`
	ExcludeRegexes []string               `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the trace's endpoints. Traces with endpoints matching regexes will be excluded"`
	Splits         []string               `json:"splits" jsonschema:"description=The splits to apply to trace timeseries data. e.g. if you want to split the trace timeseries data by service name you would set splits as ['service.name']. This is useful for seeing a breakdown of the trace timeseries data by an attribute"`
	Functions      []model.MetricFunction `json:"functions" jsonschema:"description=The functions to apply to the traces. Available functions are monotonicDifference which will calculate the difference between the current and previous value of the metric (negative values will be set to 0) and valueDifference which will calculate the difference between the current and previous value of the metric or MathExpression e.g. a / 60"`
	Aggregate      string                 `json:"aggregate" jsonschema:"required,description=The aggregation to apply to the metrics. Possible values are: count / p50 / p90 / p95 / p99 / totalSize / responseSize / requestSize. The aggregation will be applied to every datapoint bucket. For example if the bucket size is 1 minute and the aggregation is count then the count of all datapoints in a minute will be returned"`
	Environments   []string               `json:"environments" jsonschema:"description=The environments to get traces from. If empty traces from all environments will be returned"`
}

func GetTraceMetricHandler(arguments GetTraceMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetTraceMetricRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		ServiceNames:   arguments.ServiceNames,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Splits:         arguments.Splits,
		Functions:      arguments.Functions,
		Aggregate:      model.Aggregation(arguments.Aggregate),
		Environments:   arguments.Environments,
	}

	body, err := getTraceMetricMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting trace metric: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getTraceMetricMetoroCall(request model.GetTraceMetricRequest) ([]byte, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "traceMetric", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return resp, nil
}
