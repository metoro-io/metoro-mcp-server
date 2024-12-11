package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

func getAlertFiresMetoroCall(alertId string, startTime, endTime int64) ([]byte, error) {
	return MakeMetoroAPIRequest("GET", fmt.Sprintf("alertFires?alertId=%s&startTime=%d&endTime=%d", alertId, startTime, endTime), nil)
}

func getAlertFiresHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	alertId, ok := arguments["alertId"].(string)
	if !ok || alertId == "" {
		return nil, fmt.Errorf("alertId is required")
	}

	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	startTime := fiveMinsAgo.Unix()
	endTime := now.Unix()

	body, err := getAlertFiresMetoroCall(alertId, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error getting alert fires: %v", err)
	}
	return mcp.NewToolResultText(string(body)), nil
}
