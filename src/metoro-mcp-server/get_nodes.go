package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

func getNodesMetoroCall(request GetAllNodesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling nodes request: %v", err)
	}
	return MakeMetoroAPIRequest("POST", "infrastructure/nodes", bytes.NewBuffer(requestBody))
}

func getNodesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetAllNodesRequest{
		StartTime: fiveMinsAgo.Unix(),
		EndTime:   now.Unix(),
	}

	if filtersStr, ok := arguments["filters"].(string); ok && filtersStr != "" {
		var filters map[string][]string
		if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
			return nil, fmt.Errorf("error parsing filters JSON: %v", err)
		}
		request.Filters = filters
	}

	if excludeFiltersStr, ok := arguments["excludeFilters"].(string); ok && excludeFiltersStr != "" {
		var excludeFilters map[string][]string
		if err := json.Unmarshal([]byte(excludeFiltersStr), &excludeFilters); err != nil {
			return nil, fmt.Errorf("error parsing excludeFilters JSON: %v", err)
		}
		request.ExcludeFilters = excludeFilters
	}

	if splitsStr, ok := arguments["splits"].(string); ok && splitsStr != "" {
		var splits []string
		if err := json.Unmarshal([]byte(splitsStr), &splits); err != nil {
			return nil, fmt.Errorf("error parsing splits JSON: %v", err)
		}
		request.Splits = splits
	}

	if environmentsStr, ok := arguments["environments"].(string); ok && environmentsStr != "" {
		var environments []string
		if err := json.Unmarshal([]byte(environmentsStr), &environments); err != nil {
			return nil, fmt.Errorf("error parsing environments JSON: %v", err)
		}
		request.Environments = environments
	}

	body, err := getNodesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting nodes: %v", err)
	}
	return mcp.NewToolResultText(string(body)), nil
}
