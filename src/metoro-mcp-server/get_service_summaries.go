package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"time"
)

type GetServiceSummariesHandlerArgs struct {
	Namespaces   string   `json:"namespace" jsonschema:"description=The namespace to get service summaries for. If empty, all namespaces will be used."`
	Environments []string `json:"environments" jsonschema:"description=The environments to get service summaries for. If empty, all environments will be used."`
}

func getServiceSummariesHandler(arguments GetServiceSummariesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetServiceSummariesRequest{
		StartTime:    fiveMinsAgo.Unix(),
		EndTime:      now.Unix(),
		Namespace:    arguments.Namespaces,
		Environments: arguments.Environments,
	}

	body, err := getServiceSummariesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting service summaries: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getServiceSummariesMetoroCall(request GetServiceSummariesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling service summaries request: %v", err)
	}
	return MakeMetoroAPIRequest("POST", "serviceSummaries", bytes.NewBuffer(requestBody))
}
