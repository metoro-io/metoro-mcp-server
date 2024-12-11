package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"time"
)

// TODO: If we figure out how to input start and end times we can directly use GetLogsRequest struct.
type GetLogsHandlerArgs struct {
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the logs. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The filters to exclude from the logs. e.g., '{\"service.name\": [\"/k8s/namespaceX/serviceX\"]}' should exclude logs for serviceX in namespaceX"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the log's messages. Logs with message matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonshcema:"description=The regexes to exclude from the log's messages. Logs with message matching regexes will be excluded"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, logs will be returned in ascending order, otherwise in descending order"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get logs from. If empty, logs from all environments will be returned"`
}

func getLogsHandler(arguments GetLogsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetLogsRequest{
		StartTime:      fiveMinsAgo.Unix(),
		EndTime:        now.Unix(),
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Ascending:      arguments.Ascending,
		Environments:   arguments.Environments,
	}

	resp, err := getLogsMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting logs: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func getLogsMetoroCall(request GetLogsRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling logs request: %v", err)
	}
	return MakeMetoroAPIRequest("POST", "logs", bytes.NewBuffer(requestBody))
}
