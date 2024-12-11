package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNamespacesHandlerArgs struct{}

func GetNamespacesHandler(arguments GetNamespacesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getNamespacesMetoroCall()
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNamespacesMetoroCall() ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "namespaces", nil)
}
