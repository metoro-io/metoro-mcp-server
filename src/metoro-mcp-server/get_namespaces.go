package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetNamespacesHandlerArgs struct{}

func getNamespacesHandler(arguments GetNamespacesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getNamespacesMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNamespacesMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "namespaces", nil)
}
