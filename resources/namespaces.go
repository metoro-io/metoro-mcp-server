package resources

import (
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-mcp-server/utils"
)

func NamespacesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	response, err := utils.MakeMetoroAPIRequest("GET", "namespaces", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://namespaces", string(response), "text/plain")), nil
}
