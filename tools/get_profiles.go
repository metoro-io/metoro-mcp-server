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

type GetProfileHandlerArgs struct {
	TimeConfig     utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get the profiles data. e.g. if you want to get profiles for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	ServiceName    string           `json:"serviceName" jsonschema:"required,description=The name of the service to get profiles for"`
	ContainerNames []string         `json:"containerNames" jsonschema:"description=The container names to get profiles for"`
}

func GetProfilesHandler(ctx context.Context, arguments GetProfileHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetProfileRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		ServiceName:    arguments.ServiceName,
		ContainerNames: arguments.ContainerNames,
	}

	body, err := getProfilesMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting profiles: %v", err)
	}

	if len(body) > 200000 {
		return nil, fmt.Errorf("response too large, please refine your query to get a smaller response")
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getProfilesMetoroCall(ctx context.Context, request model.GetProfileRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling profiles request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "profiles", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
