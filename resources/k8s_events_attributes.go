package resources

import (
	"context"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func K8sEventsAttributesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	ctx := context.Background()
	resp, err := utils.MakeMetoroAPIRequest("GET", "k8s/events/summaryAttributes", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://k8sEventAttributes", string(resp), "text/plain")), nil
}
