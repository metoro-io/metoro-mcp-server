package resources

import (
	"context"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TraceAttributesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	ctx := context.Background()
	resp, err := utils.MakeMetoroAPIRequest("GET", "tracesSummaryAttributes", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://traceAttributes", string(resp), "text/plain")), nil
}
