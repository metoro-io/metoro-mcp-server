package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sEventsHandlerArgs struct {
	TimeConfig     utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get k8s events for. e.g. if you want to get k8s events for the last 6 hours you would set time_period=6 and time_window=Hours. You can also set an absoulute time range by setting start_time and end_time"`
	Filters        []model.Filter   `json:"filters" jsonschema:"description=Filters to apply to the events. Only the events that match these filters will be returned. Get the possible filter keys from the get_k8s_events_attributes tool and possible filter values from the get_k8s_event_attribute_values tool (for a filter key)"`
	ExcludeFilters []model.Filter   `json:"exclude_filters" jsonschema:"description=Filters to exclude the events. Events matching the exclude filters will not be returned. Get the possible exclude filter keys from the get_k8s_events_attributes tool and possible exclude filter values from the get_k8s_event_attribute_values tool (for a key)"`
	Regexes        []string         `json:"regexes" jsonschema:"description=Regexes to apply to the event messages. Only the events with messages that match these regexes will be returned. Regexes are ORed together. For example if you want to get events with messages that contain 'error' or 'warning' you would set the regexes as ['error' 'warning']"`
	ExcludeRegexes []string         `json:"exclude_regexes" jsonschema:"description=Regexes to exclude the events. Events with messages that match these regexes will not be returned. Exclude regexes are ANDed together."`
	Environments   []string         `json:"environments" jsonschema:"description=Environments/Clusters to get events for"`
}

func GetK8sEventsHandler(ctx context.Context, arguments GetK8sEventsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	// Convert Filter slice to map format for internal API
	filters := model.FiltersToMap(arguments.Filters)
	excludeFilters := model.FiltersToMap(arguments.ExcludeFilters)

	request := model.GetK8sEventsRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        filters,
		ExcludeFilters: excludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Environments:   arguments.Environments,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/events", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
