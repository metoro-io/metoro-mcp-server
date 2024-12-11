package resources

import (
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
)

func EnvironmentResourceHandler() (*mcpgolang.ResourceResponse, error) {
	response, err := utils.MakeMetoroAPIRequest("GET", "environments", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://environments", string(response), "text/plain")), nil

}
