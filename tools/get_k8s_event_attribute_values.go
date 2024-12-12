package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sEventAttributeValueHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get the possible values of K8 event attributes values. e.g. if you want to see the possible values for the attributes in the last 5 minutes, you would set time_period=5 and time_window=Minutes"`
	Attribute      string              `json:"attribute" jsonschema:"required,description=The attribute to get values for"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the events. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the events. Events matching the exclude filters will not be returned. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the event messages. Events with messages matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the event messages. Events with messages matching regexes will be excluded"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get events from. If empty, events from all environments will be returned"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, events will be returned in ascending order, otherwise in descending order"`
	PrevEndTime    *float64            `json:"prevEndTime" jsonschema:"description=The end time of the previous request"`
}

func GetK8sEventAttributeValuesForIndividualAttributeHandler(arguments GetK8sEventAttributeValueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime := utils.CalculateTimeRange(arguments.TimeConfig)
	request := model.GetSingleK8sEventSummaryRequest{
		GetK8sEventsRequest: model.GetK8sEventsRequest{
			StartTime:      startTime,
			EndTime:        endTime,
			Filters:        arguments.Filters,
			ExcludeFilters: arguments.ExcludeFilters,
			Regexes:        arguments.Regexes,
			ExcludeRegexes: arguments.ExcludeRegexes,
			Environments:   arguments.Environments,
			Ascending:      arguments.Ascending,
		},
		Attribute: arguments.Attribute,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/events/summaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
