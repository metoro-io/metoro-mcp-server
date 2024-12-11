package resources

import (
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-mcp-server/utils"
)

func K8sEventsAttributesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "k8s/events/summaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://k8sEventAttributes", string(resp), "text/plain")), nil
}
