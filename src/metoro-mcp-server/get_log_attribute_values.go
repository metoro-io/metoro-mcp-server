package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

func getLogAttributeValuesForIndividualAttributeHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetSingleLogSummaryRequest{
		LogSummaryRequest: LogSummaryRequest{
			StartTime: fiveMinsAgo.Unix(),
			EndTime:   now.Unix(),
		},
	}

	// Attribute is required
	if attribute, ok := arguments["attribute"].(string); ok && attribute != "" {
		request.Attribute = attribute
	} else {
		return nil, fmt.Errorf("attribute is required")
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

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := MakeMetoroAPIRequest("POST", "logsSummaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcp.NewToolResultText(string(resp)), nil
}
