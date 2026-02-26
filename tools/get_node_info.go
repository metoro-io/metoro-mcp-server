package tools

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNodeInfoHandlerArgs struct {
	TimeConfig utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get node information for. End time is used as the point in time for the node snapshot"`
	NodeName   string           `json:"nodeName" jsonschema:"required,description=The name of the node to get the YAML/information for"`
}

func GetNodeInfoHandler(ctx context.Context, arguments GetNodeInfoHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getNodeInfoMetoroCall(ctx, arguments)
	if err != nil {
		return nil, fmt.Errorf("error getting node info: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNodeInfoMetoroCall(ctx context.Context, arguments GetNodeInfoHandlerArgs) ([]byte, error) {
	path, err := buildGetNodeInfoPath(arguments)
	if err != nil {
		return nil, err
	}

	return utils.MakeMetoroAPIRequest("GET", path, nil, utils.GetAPIRequirementsFromRequest(ctx))
}

func buildGetNodeInfoPath(arguments GetNodeInfoHandlerArgs) (string, error) {
	_, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return "", fmt.Errorf("error calculating time range: %v", err)
	}

	query := url.Values{}
	query.Set("nodeName", arguments.NodeName)
	// Set both query params to support old and new backend contracts.
	query.Set("time", strconv.FormatInt(endTime, 10))
	query.Set("startTime", strconv.FormatInt(endTime, 10))

	return "infrastructure/node?" + query.Encode(), nil
}
