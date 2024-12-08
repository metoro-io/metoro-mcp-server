package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func getServicesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	body, err := getServicesMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting services: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf(" %s", string(body))), nil
}

func getServicesMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "services", nil)
}
