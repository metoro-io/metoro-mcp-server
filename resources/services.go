package resources

import (
	"context"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func ServicesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	ctx := context.Background()
	response, err := utils.MakeMetoroAPIRequest("GET", "services", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://services", string(response), "text/plain")), nil
}
