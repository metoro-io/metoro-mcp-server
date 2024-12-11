package tools

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-mcp-server/utils"
	"time"
)

type GetAlertFiresHandlerArgs struct {
	AlertId string `json:"alertId" jsonschema:"required,description=The alert ID to get fires for"`
}

func GetAlertFiresHandler(arguments GetAlertFiresHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)

	body, err := getAlertFiresMetoroCall(arguments.AlertId, fiveMinsAgo.Unix(), now.Unix())
	if err != nil {
		return nil, fmt.Errorf("error getting alert fires: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getAlertFiresMetoroCall(alertId string, startTime, endTime int64) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", fmt.Sprintf("alertFires?alertId=%s&startTime=%d&endTime=%d", alertId, startTime, endTime), nil)
}
