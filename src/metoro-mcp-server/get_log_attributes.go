package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func getLogAttributesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	resp, err := MakeMetoroAPIRequest("GET", "logsSummaryAttributes", nil)
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcp.NewToolResultText(string(resp)), nil
}
