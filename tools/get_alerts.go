package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetAlertHandlerArgs struct{}

func GetAlertsHandler(ctx context.Context, arguments GetAlertHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getAlertsMetoroCall(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting alerts: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getAlertsMetoroCall(ctx context.Context) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "searchAlerts", nil, utils.GetAPIRequirementsFromRequest(ctx))
}
