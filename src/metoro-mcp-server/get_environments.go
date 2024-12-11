package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetEnvironmentHandlerArgs struct{}

func getEnvironmentsHandler(arguments GetEnvironmentHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getEnvironmentsMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting environments: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getEnvironmentsMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "environments", nil)
}
