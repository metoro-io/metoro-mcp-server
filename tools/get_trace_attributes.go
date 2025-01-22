package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetTraceAttributesHandlerArgs struct{}

func GetTraceAttributesHandler(ctx context.Context, arguments GetTraceAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "tracesSummaryAttributes", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
