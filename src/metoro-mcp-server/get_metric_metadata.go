package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

type MetricMetadata struct {
	Type        string `json:"type"`
	Unit        string `json:"metricUnit"`
	Description string `json:"metricDescription"`
}

func getMetricMetadata(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	metricName := arguments["name"].(string)
	if metricName == "" {
		return nil, fmt.Errorf("metricName is required")
	}
	response, err := getMetricMetadataMetoroCall(metricName)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro get metric metadata api: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("%s", string(response))), nil
}

func getMetricMetadataMetoroCall(metricName string) ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "metric/metadata?name="+metricName, nil)
}
