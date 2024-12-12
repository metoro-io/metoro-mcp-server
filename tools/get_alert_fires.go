package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetAlertFiresHandlerArgs struct {
	TimeConfig utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get alert fires for. e.g. if you want to get fires for the last 5 minutes, you would set time_period=5 and time_window=Minutes"`
	AlertId    string           `json:"alert_id" jsonschema:"required,description=The ID of the alert to get fires for"`
}

func GetAlertFiresHandler(arguments GetAlertFiresHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	body, err := getAlertFiresMetoroCall(arguments.AlertId, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error getting alert fires: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getAlertFiresMetoroCall(alertId string, startTime, endTime int64) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", fmt.Sprintf("alertFires?alertId=%s&startTime=%d&endTime=%d", alertId, startTime, endTime), nil)
}
