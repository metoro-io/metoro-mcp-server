package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func getNamespacesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	body, err := getNamespacesMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s", string(body))), nil
}

func getNamespacesMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "namespaces", nil)
}
