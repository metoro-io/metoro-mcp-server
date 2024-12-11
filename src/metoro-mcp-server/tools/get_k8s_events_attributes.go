package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
)

type GetK8sEventsAttributesHandlerArgs struct{}

func GetK8sEventsAttributesHandler(arguments GetK8sEventsAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "k8s/events/summaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
