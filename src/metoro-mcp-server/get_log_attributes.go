package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetLogAttributesHandlerArgs struct{}

func getLogAttributesHandler(arguments GetLogAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := MakeMetoroAPIRequest("GET", "logsSummaryAttributes", nil)
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
