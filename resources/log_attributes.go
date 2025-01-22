package resources

import (
	"context"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func LogAttributesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	ctx := context.Background()
	resp, err := utils.MakeMetoroAPIRequest("GET", "logsSummaryAttributes", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://logAttributes", string(resp), "text/plain")), nil
}
