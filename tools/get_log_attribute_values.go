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

type GetLogAttributeValuesHandlerArgs struct {
	TimeConfig     utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to use while getting the possible values of log attributes. e.g. if you want to get values for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Attribute      string           `json:"attribute" jsonschema:"required,description=The attribute key to get the possible values for"`
	Filters        []model.Filter   `json:"filters" jsonschema:"description=The filters to apply before getting the possible values. For example if you want to get the possible values for attribute key service.name where the environment is X you would set the Filters as [{key: 'environment' values: ['X']}]"`
	ExcludeFilters []model.Filter   `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate possible values an attribute can take. Log attributes matching the exclude filters will not be returned."`
	Regexes        []string         `json:"regexes" jsonschema:"description=The regexes to apply to the log messages. Only the attribute values (for a given attribute key) of logs messages that match these regexes will be returned."`
	ExcludeRegexes []string         `json:"excludeRegexes" jsonschema:"description=The exclude regexes to apply to the log messages. The attribute values (for a given attribute key) of log messages that match these regexes will not be returned."`
	Environments   []string         `json:"environments" jsonschema:"description=The environments to get possible values of a log attributes for. If empty then possible values from all environments will be returned"`
}

func GetLogAttributeValuesForIndividualAttributeHandler(ctx context.Context, arguments GetLogAttributeValuesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	// Convert Filter slice to map format for internal API
	filters := model.FiltersToMap(arguments.Filters)
	excludeFilters := model.FiltersToMap(arguments.ExcludeFilters)

	request := model.GetSingleLogSummaryRequest{
		LogSummaryRequest: model.LogSummaryRequest{
			StartTime:      startTime,
			EndTime:        endTime,
			Filters:        filters,
			ExcludeFilters: excludeFilters,
			Regexes:        arguments.Regexes,
			ExcludeRegexes: arguments.ExcludeRegexes,
			Environments:   arguments.Environments,
		},
		Attribute: arguments.Attribute,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "logsSummaryIndividualAttribute", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
