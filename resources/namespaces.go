package resources

import (
	"context"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func NamespacesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	ctx := context.Background()
	response, err := utils.MakeMetoroAPIRequest("GET", "namespaces", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://namespaces", string(response), "text/plain")), nil
}
