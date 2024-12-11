package resources

import (
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-mcp-server/utils"
)

func LogAttributesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "logsSummaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://logAttributes", string(resp), "text/plain")), nil
}
