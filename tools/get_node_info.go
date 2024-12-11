package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
	"time"
)

type GetNodeInfoHandlerArgs struct {
	NodeName string `json:"nodeName" jsonschema:"required,description=The name of the node to get info for"`
}

func GetNodeInfoHandler(arguments GetNodeInfoHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)

	body, err := getNodeInfoMetoroCall(arguments.NodeName, fiveMinsAgo.Unix())
	if err != nil {
		return nil, fmt.Errorf("error getting node info: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNodeInfoMetoroCall(nodeName string, startTime int64) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", fmt.Sprintf("infrastructure/node?nodeName=%s&startTime=%d", nodeName, startTime), nil)
}
