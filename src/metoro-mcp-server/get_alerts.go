package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func getAlertsMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "searchAlerts", nil)
}

func getAlertsHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	body, err := getAlertsMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting alerts: %v", err)
	}
	return mcp.NewToolResultText(string(body)), nil
}
