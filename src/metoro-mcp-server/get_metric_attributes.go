package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

type MetricAttributesRequest struct {
	StartTime        int64               `json:"startTime"`
	EndTime          int64               `json:"endTime"`
	MetricName       string              `json:"metricName"`
	FilterAttributes map[string][]string `json:"filterAttributes"`
}

func getMetricAttributesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := MetricAttributesRequest{
		StartTime: fiveMinsAgo.Unix(),
		EndTime:   now.Unix(),
	}

	if metricName, ok := arguments["metricName"].(string); ok && metricName != "" {
		request.MetricName = metricName
	}

	if filterAttributesStr, ok := arguments["filterAttributes"].(string); ok && filterAttributesStr != "" {
		var filterAttributes map[string][]string
		if err := json.Unmarshal([]byte(filterAttributesStr), &filterAttributes); err != nil {
			return nil, fmt.Errorf("error parsing filterAttributes JSON: %v", err)
		}
		request.FilterAttributes = filterAttributes
	}

	response, err := getMetricAttributesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s", string(response))), nil
}

func getMetricAttributesMetoroCall(request MetricAttributesRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return MakeMetoroAPIRequest("POST", "metricAttributes", bytes.NewBuffer(jsonData))
}
