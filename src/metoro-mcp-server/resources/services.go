package resources

import (
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
)

func ServicesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	response, err := utils.MakeMetoroAPIRequest("GET", "services", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://services", string(response), "text/plain")), nil
}
