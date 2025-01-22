package tools

import (
	"context"
	"fmt"
	"time"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNodeInfoHandlerArgs struct {
	NodeName string `json:"nodeName" jsonschema:"required,description=The name of the node to get the YAML/information for"`
}

func GetNodeInfoHandler(ctx context.Context, arguments GetNodeInfoHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-10 * time.Minute)

	body, err := getNodeInfoMetoroCall(ctx, arguments.NodeName, fiveMinsAgo.Unix())
	if err != nil {
		return nil, fmt.Errorf("error getting node info: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNodeInfoMetoroCall(ctx context.Context, nodeName string, startTime int64) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", fmt.Sprintf("infrastructure/node?nodeName=%s&startTime=%d", nodeName, startTime), nil, utils.GetAPIRequirementsFromRequest(ctx))
}
