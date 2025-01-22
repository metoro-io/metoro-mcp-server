package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetLogAttributesHandlerArgs struct{}

func GetLogAttributesHandler(ctx context.Context, arguments GetLogAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "logsSummaryAttributes", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
