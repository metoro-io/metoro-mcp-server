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

type GetServiceSummariesHandlerArgs struct {
	TimeConfig   utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get service summaries for. e.g. if you want to get summaries for the last 5 minutes you would set time_period=5 and time_window=Minutes. Try to use a time period 1 hour or less. You can also set an absoulute time range by setting start_time and end_time"`
	Namespaces   string           `json:"namespace" jsonschema:"description=The namespace to get service summaries for. If empty all namespaces will be used."`
	Environments []string         `json:"environments" jsonschema:"description=The environments to get service summaries for. If empty all environments will be used."`
}

func GetServiceSummariesHandler(ctx context.Context, arguments GetServiceSummariesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetServiceSummariesRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		Namespace:    arguments.Namespaces,
		Environments: arguments.Environments,
	}

	body, err := getServiceSummariesMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting service summaries: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getServiceSummariesMetoroCall(ctx context.Context, request model.GetServiceSummariesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling service summaries request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "serviceSummaries", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
