package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

type Aggregation string
type MetricFunction struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}

func getMetricHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetMetricRequest{
		StartTime: fiveMinsAgo.Unix(),
		EndTime:   now.Unix(),
	}

	if metricName, ok := arguments["metricName"].(string); ok && metricName != "" {
		request.MetricName = metricName
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

	if aggregation, ok := arguments["aggregation"].(string); ok && aggregation != "" {
		request.Aggregation = Aggregation(aggregation)
	}

	if functionsStr, ok := arguments["functions"].(string); ok && functionsStr != "" {
		var functions []MetricFunction
		if err := json.Unmarshal([]byte(functionsStr), &functions); err != nil {
			return nil, fmt.Errorf("error parsing functions JSON: %v", err)
		}
		request.Functions = functions
	}

	if limitResults, ok := arguments["limitResults"].(bool); ok {
		request.LimitResults = limitResults
	}

	if bucketSize, ok := arguments["bucketSize"].(float64); ok {
		request.BucketSize = int64(bucketSize)
	}

	body, err := getMetricMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting metric: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s", string(body))), nil
}

func getMetricMetoroCall(request GetMetricRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling metric request: %v", err)
	}
	return MakeMetoroAPIRequest("POST", "metric", bytes.NewBuffer(requestBody))
}
