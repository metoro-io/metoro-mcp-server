package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"time"
)

type GetK8sEventsHandlerArgs struct {
	Filters        map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the events"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=Filters to exclude from the events"`
	Regexes        []string            `json:"regexes" jsonschema:"description=Regexes to apply to the event messages"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=Regexes to exclude from the event messages"`
	Environments   []string            `json:"environments" jsonschema:"description=Environments to get events from"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, events will be returned in ascending order, otherwise in descending order"`
	PrevEndTime    *float64            `json:"prevEndTime" jsonschema:"description=The end time of the previous request"`
}

func getK8sEventsHandler(arguments GetK8sEventsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	sixHoursAgo := now.Add(-6 * time.Hour)
	request := GetK8sEventsRequest{
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
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := MakeMetoroAPIRequest("POST", "k8s/events", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
