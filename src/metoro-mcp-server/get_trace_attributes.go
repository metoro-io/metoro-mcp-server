package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetTraceAttributesHandlerArgs struct{}

func getTraceAttributesHandler(arguments GetTraceAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := MakeMetoroAPIRequest("GET", "tracesSummaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
