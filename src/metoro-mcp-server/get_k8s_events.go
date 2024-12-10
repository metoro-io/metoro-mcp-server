package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

func getK8sEventsHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	sixHoursAgo := now.Add(-6 * time.Hour)
	request := GetK8sEventsRequest{
		StartTime: sixHoursAgo.Unix(),
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

	if regexesStr, ok := arguments["regexes"].(string); ok && regexesStr != "" {
		var regexes []string
		if err := json.Unmarshal([]byte(regexesStr), &regexes); err != nil {
			return nil, fmt.Errorf("error parsing regexes JSON: %v", err)
		}
		request.Regexes = regexes
	}

	if excludeRegexesStr, ok := arguments["excludeRegexes"].(string); ok && excludeRegexesStr != "" {
		var excludeRegexes []string
		if err := json.Unmarshal([]byte(excludeRegexesStr), &excludeRegexes); err != nil {
			return nil, fmt.Errorf("error parsing excludeRegexes JSON: %v", err)
		}
		request.ExcludeRegexes = excludeRegexes
	}

	if environmentsStr, ok := arguments["environments"].(string); ok && environmentsStr != "" {
		var environments []string
		if err := json.Unmarshal([]byte(environmentsStr), &environments); err != nil {
			return nil, fmt.Errorf("error parsing environments JSON: %v", err)
		}
		request.Environments = environments
	}

	if ascending, ok := arguments["ascending"].(bool); ok {
		request.Ascending = ascending
	}

	if prevEndTimeFloat, ok := arguments["prevEndTime"].(float64); ok {
		prevEndTime := int64(prevEndTimeFloat)
		request.PrevEndTime = &prevEndTime
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := MakeMetoroAPIRequest("POST", "k8s/events", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcp.NewToolResultText(string(resp)), nil
}
