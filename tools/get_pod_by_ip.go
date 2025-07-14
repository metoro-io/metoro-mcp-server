package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetResourcesByIpHandlerArgs struct {
	Ip          string           `json:"ip" jsonschema:"required,description=IP address to search for"`
	TimeConfig  utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get the resources for. e.g. if you want the get the resources for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	Environment string           `json:"environment" jsonschema:"required,description=Environment to filter the resources by"`
}

type GetResourcesByIpRequest struct {
	Ip          string `json:"ip"`
	StartTime   int64  `json:"startTime"`
	EndTime     int64  `json:"endTime"`
	Environment string `json:"environment"`
}

type GetResourcesByIpResponse struct {
	Resources []ResourcesByIpData `json:"resources"`
}

type ResourcesByIpData struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Ip          string `json:"ip"`
	NodeName    string `json:"nodeName"`
	Status      string `json:"status"`
	Environment string `json:"environment"`
}

func GetResourcesByIpHandler(ctx context.Context, arguments GetResourcesByIpHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetResourcesByIpRequest{
		Ip:          arguments.Ip,
		StartTime:   startTime,
		EndTime:     endTime,
		Environment: arguments.Environment,
	}

	resp, err := getPodByIpMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting pod by IP: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func getPodByIpMetoroCall(ctx context.Context, request GetResourcesByIpRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling pod by IP request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "k8s/resources/byIp", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
