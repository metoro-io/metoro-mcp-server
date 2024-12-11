package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetServicesHandlerArgs struct{}

func getServicesHandler(arguments GetServicesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getServicesMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting services: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getServicesMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "services", nil)
}
