package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

type FuzzyMetricsRequest struct {
	MetricFuzzyMatch string   `json:"metricFuzzyMatch"`
	Environments     []string `json:"environments"`
	StartTime        int64    `json:"startTime"`
	EndTime          int64    `json:"endTime"`
}

func getMetricNamesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := FuzzyMetricsRequest{
		StartTime: fiveMinsAgo.Unix(),
		EndTime:   now.Unix(),
	}

	//if fuzzyMatch, ok := arguments["metricFuzzyMatch"].(string); ok {
	//	request.MetricFuzzyMatch = fuzzyMatch
	//}

	if environmentsStr, ok := arguments["environments"].(string); ok && environmentsStr != "" {
		var environments []string
		if err := json.Unmarshal([]byte(environmentsStr), &environments); err != nil {
			return nil, fmt.Errorf("error parsing environments JSON: %v", err)
		}
		request.Environments = environments
	}
	request.MetricFuzzyMatch = ""

	response, err := getMetricNamesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s", string(response))), nil
}

func getMetricNamesMetoroCall(request FuzzyMetricsRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return MakeMetoroAPIRequest("POST", "/fuzzyMetricsNames", bytes.NewBuffer(jsonData))
}
