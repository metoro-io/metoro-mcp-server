package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

type GetServiceSummariesRequest struct {
	// Required: Start time of when to get the service summaries in seconds
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds
	EndTime int64 `json:"endTime"`
	// If empty, all services across all environments will be returned
	Environments []string `json:"environments"`
	// Required: The namespace of the services to get summaries for. If empty, return services from all namespaces
	Namespace string `json:"namespace"`
}

func getServiceSummariesMetoroCall(request GetServiceSummariesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling service summaries request: %v", err)
	}
	return MakeMetoroAPIRequest("POST", "serviceSummaries", bytes.NewBuffer(requestBody))
}

func getServiceSummariesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetServiceSummariesRequest{
		StartTime: fiveMinsAgo.Unix(),
		EndTime:   now.Unix(),
	}

	if namespace, ok := arguments["namespace"].(string); ok {
		request.Namespace = namespace
	}

	if environmentsStr, ok := arguments["environments"].(string); ok && environmentsStr != "" {
		var environments []string
		if err := json.Unmarshal([]byte(environmentsStr), &environments); err != nil {
			return nil, fmt.Errorf("error parsing environments JSON: %v", err)
		}
		request.Environments = environments
	}

	body, err := getServiceSummariesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting service summaries: %v", err)
	}
	return mcp.NewToolResultText(string(body)), nil
}
