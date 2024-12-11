package resources

import (
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-mcp-server/utils"
)

func ServicesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	response, err := utils.MakeMetoroAPIRequest("GET", "services", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://services", string(response), "text/plain")), nil
}
