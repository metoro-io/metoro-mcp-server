package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetLogsHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get the logs for. e.g. if you want the get the logs for the last 5 minutes, you would set time_period=5 and time_window=Minutes"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the logs. it is a map of filter keys to array values where array values are ORed.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The filters to exclude from the logs. e.g., '{\"service.name\": [\"/k8s/namespaceX/serviceX\"]}' should exclude logs for serviceX in namespaceX"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the log's messages. Logs with message matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonshcema:"description=The regexes to exclude from the log's messages. Logs with message matching regexes will be excluded"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, logs will be returned in ascending order, otherwise in descending order"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get logs from. If empty, logs from all environments will be returned"`
}

func GetLogsHandler(arguments GetLogsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime := utils.CalculateTimeRange(arguments.TimeConfig)

	request := model.GetLogsRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Ascending:      arguments.Ascending,
		Environments:   arguments.Environments,
	}

	resp, err := getLogsMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting logs: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func getLogsMetoroCall(request model.GetLogsRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling logs request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "logs", bytes.NewBuffer(requestBody))
}
