package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/model"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
	"time"
)

type GetTraceMetricHandlerArgs struct {
	ServiceNames   []string               `json:"serviceNames" jsonschema:"description=Service names to return traces for"`
	Filters        map[string][]string    `json:"filters" jsonschema:"description=The filters to apply to the traces. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string    `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the traces. Traces matching the exclude traces will not be returned. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	Regexes        []string               `json:"regexes" jsonschema:"description=The regexes to apply to the trace's endpoints. Traces with endpoints matching regexes will be returned"`
	ExcludeRegexes []string               `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the trace's endpoints. Traces with endpoints matching regexes will be excluded"`
	Splits         []string               `json:"splits" jsonschema:"description=The splits to apply to the metric. Metrics will be split by the given keys"`
	Functions      []model.MetricFunction `json:"functions" jsonschema:"description=The functions to apply to the metric"`
	Aggregate      string                 `json:"aggregate" jsonschema:"description=The aggregation to apply to the metric. e.g. sum, avg, min, max, count"`
	Environments   []string               `json:"environments" jsonschema:"description=The environments to get traces from. If empty, traces from all environments will be returned"`
}

func GetTraceMetricHandler(arguments GetTraceMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := model.GetTraceMetricRequest{
		StartTime:      fiveMinsAgo.Unix(),
		EndTime:        now.Unix(),
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
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
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
