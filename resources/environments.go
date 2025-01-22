package resources

import (
	"context"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func EnvironmentResourceHandler() (*mcpgolang.ResourceResponse, error) {
	ctx := context.Background()
	response, err := utils.MakeMetoroAPIRequest("GET", "environments", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://environments", string(response), "text/plain")), nil

}
