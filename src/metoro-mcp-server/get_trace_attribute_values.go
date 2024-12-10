package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

type GetSingleTraceSummaryRequest struct {
	TracesSummaryRequest
	// The attribute to get the summary for
	Attribute string `json:"attribute"`
}

type TracesSummaryRequest struct {
	// Required: Start time of when to get the service summaries in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds since epoch
	EndTime int64 `json:"endTime"`

	// The filters to apply to the trace summary, so for example, if you want to get traces for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`
	// ExcludeFilters are used to exclude traces based on a filter
	ExcludeFilters map[string][]string `json:"excludeFilters"`

	// Regexes are used to filter traces based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter traces based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`

	// Optional: The name of the service to get the trace metrics for
	// Acts as an additional filter
	ServiceNames []string `json:"serviceNames"`

	// Environments is the environments to get the traces for. If empty, all environments will be included
	Environments []string `json:"environments"`
}

func getTraceAttributeValuesForIndividualAttributeHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetSingleTraceSummaryRequest{
		TracesSummaryRequest: TracesSummaryRequest{
			StartTime: fiveMinsAgo.Unix(),
			EndTime:   now.Unix(),
		},
	}

	// Required attribute parameter
	if attribute, ok := arguments["attribute"].(string); ok && attribute != "" {
		request.Attribute = attribute
	} else {
		return nil, fmt.Errorf("attribute is required")
	}

	if serviceNamesStr, ok := arguments["serviceNames"].(string); ok && serviceNamesStr != "" {
		var serviceNames []string
		if err := json.Unmarshal([]byte(serviceNamesStr), &serviceNames); err != nil {
			return nil, fmt.Errorf("error parsing serviceNames JSON: %v", err)
		}
		request.ServiceNames = serviceNames
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

	resp, err := MakeMetoroAPIRequest("POST", "tracesSummaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcp.NewToolResultText(string(resp)), nil
}
