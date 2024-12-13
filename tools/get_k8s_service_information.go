package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sServiceInformationHandlerArgs struct {
	TimeConfig   utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time to get state of the YAML file. e.g. if you want to see the state of the service 5 minutes ago you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	ServiceName  string           `json:"serviceName" jsonschema:"required,description=The name of the service to get YAML file for."`
	Environments []string         `json:"environments" jsonschema:"description=The environments to get service YAML for. If empty all environments will be used."`
}

func GetK8sServiceInformationHandler(arguments GetK8sServiceInformationHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetPodsRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		ServiceName:  arguments.ServiceName,
		Environments: arguments.Environments,
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/summary", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
