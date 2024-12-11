package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"time"
)

type GetNodesHandlerArgs struct {
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the nodes"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The filters to exclude from the nodes"`
	Splits         []string            `json:"splits" jsonschema:"description=The splits to apply to the nodes"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get nodes from"`
}

func getNodesHandler(arguments GetNodesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetAllNodesRequest{
		StartTime:      fiveMinsAgo.Unix(),
		EndTime:        now.Unix(),
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Splits:         arguments.Splits,
		Environments:   arguments.Environments,
	}
	body, err := getNodesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting nodes: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNodesMetoroCall(request GetAllNodesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling nodes request: %v", err)
	}
	return MakeMetoroAPIRequest("POST", "infrastructure/nodes", bytes.NewBuffer(requestBody))
}
