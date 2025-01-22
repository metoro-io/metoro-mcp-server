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

type GetPodsHandlerArgs struct {
	TimeConfig   utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get pods for. e.g. if you want to get pods for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	ServiceName  string           `json:"serviceName" jsonschema:"description=The name of the service to get pods for. One of serviceName or nodeName is required"`
	NodeName     string           `json:"nodeName" jsonschema:"description=The name of the node to get pods for. One of serviceName or nodeName is required"`
	Environments []string         `json:"environments" jsonschema:"description=The environments to get pods for. If empty, all environments will be used."`
}

func GetPodsHandler(ctx context.Context, arguments GetPodsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	// One of serviceName or nodeName is required.
	if arguments.ServiceName == "" && arguments.NodeName == "" {
		return nil, fmt.Errorf("one of serviceName or nodeName is required")
	}

	request := model.GetPodsRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		Environments: arguments.Environments,
		ServiceName:  arguments.ServiceName,
		NodeName:     arguments.NodeName,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/pods", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
