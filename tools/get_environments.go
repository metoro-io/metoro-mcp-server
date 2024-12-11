package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-mcp-server/utils"
)

type GetEnvironmentHandlerArgs struct{}

func GetEnvironmentsHandler(arguments GetEnvironmentHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getEnvironmentsMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting environments: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getEnvironmentsMetoroCall() ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "environments", nil)
}
