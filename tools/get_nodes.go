package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNodesHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get nodes for. e.g. if you want to get nodes for the last 5 minutes, you would set time_period=5 and time_window=Minutes"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the nodes"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The filters to exclude from the nodes"`
	Splits         []string            `json:"splits" jsonschema:"description=The splits to apply to the nodes"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get nodes from"`
}

func GetNodesHandler(arguments GetNodesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime := utils.CalculateTimeRange(arguments.TimeConfig)
	request := model.GetAllNodesRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Splits:         arguments.Splits,
		Environments:   arguments.Environments,
	}
	body, err := getNodesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting nodes: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNodesMetoroCall(request model.GetAllNodesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling nodes request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "infrastructure/nodes", bytes.NewBuffer(requestBody))
}
