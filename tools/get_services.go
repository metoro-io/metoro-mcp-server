package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetServicesHandlerArgs struct{}

func GetServicesHandler(ctx context.Context, arguments GetServicesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getServicesMetoroCall(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting services: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getServicesMetoroCall(ctx context.Context) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "services", nil, utils.GetAPIRequirementsFromRequest(ctx))
}
