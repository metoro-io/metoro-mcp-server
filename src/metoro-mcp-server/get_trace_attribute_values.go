package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"time"
)

type GetTraceAttributeValuesHandlerArgs struct {
	Attribute    string   `json:"attribute" jsonschema:"required, description=The name of the attribute to get values for"`
	ServiceNames []string `json:"serviceNames" jsonschema:"description=The service names to get attribute values for"`
	//  TODO: I don't think we need these two fields for the LLM tool
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the traces. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the traces. Traces matching the exclude traces will not be returned. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the trace's endpoints. Traces with endpoints matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the trace's endpoints. Traces with endpoints matching regexes will be excluded"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get traces from. If empty, traces from all environments will be returned"`
}

func getTraceAttributeValuesForIndividualAttributeHandler(arguments GetTraceAttributeValuesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetSingleTraceSummaryRequest{
		TracesSummaryRequest: TracesSummaryRequest{
			StartTime:      fiveMinsAgo.Unix(),
			EndTime:        now.Unix(),
			ServiceNames:   arguments.ServiceNames,
			Filters:        arguments.Filters,
			ExcludeFilters: arguments.ExcludeFilters,
			Regexes:        arguments.Regexes,
			ExcludeRegexes: arguments.ExcludeRegexes,
			Environments:   arguments.Environments,
		},
		Attribute: arguments.Attribute,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := MakeMetoroAPIRequest("POST", "tracesSummaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
