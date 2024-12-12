package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sEventsHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get events for. e.g. if you want to get events for the last 6 hours, you would set time_period=6 and time_window=Hours"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the events"`
	ExcludeFilters map[string][]string `json:"exclude_filters" jsonschema:"description=Filters to exclude from the events"`
	Regexes        []string            `json:"regexes" jsonschema:"description=Regexes to apply to the event messages"`
	ExcludeRegexes []string            `json:"exclude_regexes" jsonschema:"description=Regexes to exclude from the event messages"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, events will be returned in ascending order, otherwise in descending order"`
	Environments   []string            `json:"environments" jsonschema:"description=Environments to get events from"`
}

func GetK8sEventsHandler(arguments GetK8sEventsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime := utils.CalculateTimeRange(arguments.TimeConfig)
	request := model.GetK8sEventsRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Ascending:      arguments.Ascending,
		Environments:   arguments.Environments,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/events", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
