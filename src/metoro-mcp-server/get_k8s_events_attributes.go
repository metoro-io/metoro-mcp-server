package main

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func getK8sEventsAttributesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	resp, err := MakeMetoroAPIRequest("GET", "k8s/events/summaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resp)), nil
}
