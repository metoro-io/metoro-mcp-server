package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNodeAttributesHandlerArgs struct {
	TimeConfig utils.TimeConfig `json:"timeConfig" jsonschema:"required,description=The time range to get the node attributes for"`
}

func GetNodeAttributesHandler(arguments GetNodeAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.MetricAttributesRequest{
		StartTime:        startTime,
		EndTime:          endTime,
		MetricName:       "node_info",
		FilterAttributes: map[string][]string{},
	}
	response, err := getMetricAttributesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}
