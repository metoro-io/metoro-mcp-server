package resources

import (
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
)

func TraceAttributesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "tracesSummaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://traceAttributes", string(resp), "text/plain")), nil
}
