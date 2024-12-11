package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
)

type GetAlertHandlerArgs struct{}

func GetAlertsHandler(arguments GetAlertHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getAlertsMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting alerts: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getAlertsMetoroCall() ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "searchAlerts", nil)
}
