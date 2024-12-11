package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
	"time"
)

type GetK8sEventsVolumeHandlerArgs struct {
	Filters        map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the events"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=Filters to exclude from the events"`
	Regexes        []string            `json:"regexes" jsonschema:"description=Regexes to apply to the event messages"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=Regexes to exclude from the event messages"`
	Environments   []string            `json:"environments" jsonschema:"description=Environments to get events from"`
}

func GetK8sEventsVolumeHandler(arguments GetK8sEventsVolumeHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	sixHoursAgo := now.Add(-6 * time.Hour)
	request := model.GetK8sEventMetricsRequest{
		StartTime:      sixHoursAgo.Unix(),
		EndTime:        now.Unix(),
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Environments:   arguments.Environments,
		Splits:         []string{"EventType"}, // We want the volume to be split by EventType so we can see the breakdown of Warning/Normal events.
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/events/metrics", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
