package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetK8sEventsAttributesHandlerArgs struct{}

func getK8sEventsAttributesHandler(arguments GetK8sEventsAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := MakeMetoroAPIRequest("GET", "k8s/events/summaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
