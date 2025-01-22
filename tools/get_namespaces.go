package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNamespacesHandlerArgs struct{}

func GetNamespacesHandler(ctx context.Context, arguments GetNamespacesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getNamespacesMetoroCall(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNamespacesMetoroCall(ctx context.Context) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "namespaces", nil, utils.GetAPIRequirementsFromRequest(ctx))
}
