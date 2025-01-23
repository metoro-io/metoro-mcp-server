package tools

import (
	"context"
	"fmt"
	"strings"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNamespacesHandlerArgs struct {
	Environments []string `json:"environments" jsonschema:"description=The environments (aka clusters) to get namespaces for. If empty all environments will be used."`
}

func GetNamespacesHandler(ctx context.Context, arguments GetNamespacesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getNamespacesMetoroCall(ctx, arguments)
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNamespacesMetoroCall(ctx context.Context, arguments GetNamespacesHandlerArgs) ([]byte, error) {
	endpoint := "namespaces"
	if len(arguments.Environments) > 0 {
		endpoint = fmt.Sprintf("namespaces?environments=%s", strings.Join(arguments.Environments, ","))
	}

	return utils.MakeMetoroAPIRequest("GET", endpoint, nil, utils.GetAPIRequirementsFromRequest(ctx))
}
