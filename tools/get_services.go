package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
)

type GetServicesHandlerArgs struct{}

func GetServicesHandler(arguments GetServicesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getServicesMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting services: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getServicesMetoroCall() ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "services", nil)
}
