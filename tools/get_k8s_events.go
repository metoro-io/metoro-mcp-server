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
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get k8s events for. e.g. if you want to get k8s events for the last 6 hours you would set time_period=6 and time_window=Hours. You can also set an absoulute time range by setting start_time and end_time"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the events. Only the events that match these filters will be returned. Get the possible filter keys from the get_k8s_events_attributes tool and possible filter values from the get_k8s_event_attribute_values tool (for a filter key)"`
	ExcludeFilters map[string][]string `json:"exclude_filters" jsonschema:"description=Filters to exclude the events. Events matching the exclude filters will not be returned. Get the possible exclude filter keys from the get_k8s_events_attributes tool and possible exclude filter values from the get_k8s_event_attribute_values tool (for a key)"`
	Regexes        []string            `json:"regexes" jsonschema:"description=Regexes to apply to the event messages. Only the events with messages that match these regexes will be returned. Regexes are ORed together. For example if you want to get events with messages that contain the word 'error' or 'warning' you would set the regexes as ['error' 'warning']"`
	ExcludeRegexes []string            `json:"exclude_regexes" jsonschema:"description=Regexes to exclude the events. Events with messages that match these regexes will not be returned. Exclude regexes are AND together. For example if you want to get events with messages that do not contain the word 'error' or 'warning' you would set the exclude regexes as ['error' 'warning']"`
	Environments   []string            `json:"environments" jsonschema:"description=Environments/Clusters to get events for"`
}

func GetK8sEventsHandler(arguments GetK8sEventsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetK8sEventsRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
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
