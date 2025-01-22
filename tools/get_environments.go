package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetEnvironmentHandlerArgs struct{}

func GetEnvironmentsHandler(ctx context.Context, arguments GetEnvironmentHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getEnvironmentsMetoroCall(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting environments: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getEnvironmentsMetoroCall(ctx context.Context) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "environments", nil, utils.GetAPIRequirementsFromRequest(ctx))
}
