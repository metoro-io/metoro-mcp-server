package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"time"
)

type GetK8sEventAttributeValueHandlerArgs struct {
	Attribute      string              `json:"attribute" jsonschema:"required,description=The attribute to get values for"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the events. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the events. Events matching the exclude filters will not be returned. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the event messages. Events with messages matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the event messages. Events with messages matching regexes will be excluded"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get events from. If empty, events from all environments will be returned"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, events will be returned in ascending order, otherwise in descending order"`
	PrevEndTime    *float64            `json:"prevEndTime" jsonschema:"description=The end time of the previous request"`
}

func getK8sEventAttributeValuesForIndividualAttributeHandler(arguments GetK8sEventAttributeValueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	sixHoursAgo := now.Add(-6 * time.Hour)
	request := GetSingleK8sEventSummaryRequest{
		GetK8sEventsRequest: GetK8sEventsRequest{
			StartTime:      sixHoursAgo.Unix(),
			EndTime:        now.Unix(),
			Filters:        arguments.Filters,
			ExcludeFilters: arguments.ExcludeFilters,
			Regexes:        arguments.Regexes,
			ExcludeRegexes: arguments.ExcludeRegexes,
			Environments:   arguments.Environments,
			Ascending:      arguments.Ascending,
			// TODO: Deal with the prevend time when you are dealing with the start and endtime.
			//PrevEndTime: arguments.PrevEndTime,
		},
		Attribute: arguments.Attribute,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := MakeMetoroAPIRequest("POST", "k8s/events/summaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
