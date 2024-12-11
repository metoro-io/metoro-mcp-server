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

type GetTracesHandlerArgs struct {
	ServiceNames   []string            `json:"serviceNames" jsonschema:"description=Service names to return traces for"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the traces. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the traces. Traces matching the exclude traces will not be returned. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the trace's endpoints. Traces with endpoints matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the trace's endpoints. Traces with endpoints matching regexes will be excluded"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, traces will be returned in ascending order, otherwise in descending order"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get traces from. If empty, traces from all environments will be returned"`
}

func GetTracesHandler(arguments GetTracesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := model.GetTracesRequest{
		StartTime:      fiveMinsAgo.Unix(),
		EndTime:        now.Unix(),
		ServiceNames:   arguments.ServiceNames,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		Ascending:      arguments.Ascending,
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
