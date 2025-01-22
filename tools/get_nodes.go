package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNodesHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get nodes for. e.g. if you want to get nodes for the last 5 minutes you would set time_period=5 and time_window=Minutes or if you want to get nodes for the last 2 hours you would set time_period=2 and time_window=Hours. You can also set an absoulute time range by setting start_time and end_time"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the nodes. Only the nodes that match these filters will be returned. To get possible filter keys and values use the get_node_attributes tool."`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The filters to exclude the nodes. Nodes matching the exclude filters will not be returned. To get possible exclude filter keys and values use the get_node_attributes tool."`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get nodes that belong to. If empty all environments will be used."`
}

func GetNodesHandler(ctx context.Context, arguments GetNodesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetAllNodesRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Splits:         []string{},
		Environments:   arguments.Environments,
	}
	body, err := getNodesMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting nodes: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNodesMetoroCall(ctx context.Context, request model.GetAllNodesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling nodes request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "infrastructure/nodes", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
