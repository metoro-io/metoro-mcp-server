package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetAlertHandlerArgs struct{}

func getAlertsHandler(arguments GetAlertHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getAlertsMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting alerts: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getAlertsMetoroCall() ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "searchAlerts", nil)
}
