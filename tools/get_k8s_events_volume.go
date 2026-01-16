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

type GetK8sEventsVolumeHandlerArgs struct {
	TimeConfig     utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get events volumes for. e.g. if you want to get events for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Filters        []model.Filter   `json:"filters" jsonschema:"description=Filters to apply to the events. Only the event matching these filters will be counted. Get the possible filter keys from the get_k8s_events_attributes tool and possible filter values from the get_k8s_event_attribute_values tool (for a filter key)"`
	ExcludeFilters []model.Filter   `json:"excludeFilters" jsonschema:"description=Filters to exclude the events. Events matching the exclude filters will not be counted. Get the possible exclude filter keys from the get_k8s_events_attributes tool and possible exclude filter values from the get_k8s_event_attribute_values tool (for a key)"`
	Regexes        []string         `json:"regexes" jsonschema:"description=Only the events with messages that match these regexes will be counted"`
	ExcludeRegexes []string         `json:"excludeRegexes" jsonschema:"description=Events with messages that match these regexes will not be counted"`
	Environments   []string         `json:"environments" jsonschema:"description=Environments to get events from"`
}

func GetK8sEventsVolumeHandler(ctx context.Context, arguments GetK8sEventsVolumeHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	// Convert Filter slice to map format for internal API
	filters := model.FiltersToMap(arguments.Filters)
	excludeFilters := model.FiltersToMap(arguments.ExcludeFilters)

	request := model.GetK8sEventMetricsRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        filters,
		ExcludeFilters: excludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Environments:   arguments.Environments,
		Splits:         []string{"EventType"}, // We want the volume to be split by EventType so we can see the breakdown of Warning/Normal events.
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/events/metrics", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
