package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetTracesHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get traces for. e.g. if you want to get traces for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time. Try to use a time period 24 hours or less unless its requested."`
	ServiceNames   []string            `json:"serviceNames" jsonschema:"description=Service names to return traces for"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the traces. Only the traces that match these filters will be returned. Get the possible filter keys from the get_trace_attributes tool and possible values of a filter key from the get_trace_attribute_values_for_individual_attribute tool"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the traces. Traces matching the exclude traces will not be returned. Get the possible exclude filter keys from the get_trace_attributes tool and possible value for the key from the get_trace_attribute_values_for_individual_attribute tool"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the trace's endpoints. Traces with endpoints matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the trace's endpoints. Traces with endpoints matching regexes will be excluded"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get traces from. If empty traces from all environments will be returned"`
}

func GetTracesHandler(arguments GetTracesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetTracesRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		ServiceNames:   arguments.ServiceNames,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		Environments:   arguments.Environments,
	}

	body, err := getTracesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting traces: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getTracesMetoroCall(request model.GetTracesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling traces request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "traces", bytes.NewBuffer(requestBody))
}
