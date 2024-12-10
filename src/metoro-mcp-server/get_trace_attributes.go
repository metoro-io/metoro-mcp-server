package main

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func getTraceAttributesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	resp, err := MakeMetoroAPIRequest("GET", "tracesSummaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resp)), nil
}
