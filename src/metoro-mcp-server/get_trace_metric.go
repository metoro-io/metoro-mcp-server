package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

func getTraceMetricHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetTraceMetricRequest{
		StartTime: fiveMinsAgo.Unix(),
		EndTime:   now.Unix(),
	}

	//if serviceNamesStr, ok := arguments["serviceNames"].(string); ok && serviceNamesStr != "" {
	//	var serviceNames []string
	//	if err := json.Unmarshal([]byte(serviceNamesStr), &serviceNames); err != nil {
	//		return nil, fmt.Errorf("error parsing serviceNames JSON: %v", err)
	//	}
	//	request.ServiceNames = serviceNames
	//}

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

	if splitsStr, ok := arguments["splits"].(string); ok && splitsStr != "" {
		var splits []string
		if err := json.Unmarshal([]byte(splitsStr), &splits); err != nil {
			return nil, fmt.Errorf("error parsing splits JSON: %v", err)
		}
		request.Splits = splits
	}

	if functionsStr, ok := arguments["functions"].(string); ok && functionsStr != "" {
		var functions []MetricFunction
		if err := json.Unmarshal([]byte(functionsStr), &functions); err != nil {
			return nil, fmt.Errorf("error parsing functions JSON: %v", err)
		}
		request.Functions = functions
	}

	if aggregate, ok := arguments["aggregate"].(string); ok && aggregate != "" {
		request.Aggregate = Aggregation(aggregate)
	}

	if environmentsStr, ok := arguments["environments"].(string); ok && environmentsStr != "" {
		var environments []string
		if err := json.Unmarshal([]byte(environmentsStr), &environments); err != nil {
			return nil, fmt.Errorf("error parsing environments JSON: %v", err)
		}
		request.Environments = environments
	}

	//if limitResults, ok := arguments["limitResults"].(bool); ok {
	//	request.LimitResults = limitResults
	//}

	if bucketSize, ok := arguments["bucketSize"].(float64); ok {
		request.BucketSize = int64(bucketSize)
	}

	body, err := getTraceMetricMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting trace metric: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s", string(body))), nil
}

func getTraceMetricMetoroCall(request GetTraceMetricRequest) ([]byte, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := MakeMetoroAPIRequest("POST", "traceMetric", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return resp, nil
}
