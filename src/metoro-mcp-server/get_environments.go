package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func getEnvironmentsHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	body, err := getEnvironmentsMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting environments: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s", string(body))), nil
}

func getEnvironmentsMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "environments", nil)
}
