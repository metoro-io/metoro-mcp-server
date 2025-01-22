package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sEventsAttributesHandlerArgs struct{}

func GetK8sEventsAttributesHandler(ctx context.Context, arguments GetK8sEventsAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "k8s/events/summaryAttributes", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
