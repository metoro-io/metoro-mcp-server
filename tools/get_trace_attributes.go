package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetTraceAttributesHandlerArgs struct{}

func GetTraceAttributesHandler(arguments GetTraceAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := utils.MakeMetoroAPIRequest("GET", "tracesSummaryAttributes", nil)
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
