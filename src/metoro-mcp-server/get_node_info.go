package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

func getNodeInfoMetoroCall(nodeName string, startTime int64) ([]byte, error) {
	return MakeMetoroAPIRequest("GET", fmt.Sprintf("infrastructure/node?nodeName=%s&startTime=%d", nodeName, startTime), nil)
}

func getNodeInfoHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	nodeName, ok := arguments["nodeName"].(string)
	if !ok || nodeName == "" {
		return nil, fmt.Errorf("nodeName is required")
	}

	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)

	body, err := getNodeInfoMetoroCall(nodeName, fiveMinsAgo.Unix())
	if err != nil {
		return nil, fmt.Errorf("error getting node info: %v", err)
	}
	return mcp.NewToolResultText(string(body)), nil
}
