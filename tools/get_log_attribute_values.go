package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetLogAttributeValuesHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get log attribute values for. e.g. if you want to get values for the last 5 minutes, you would set time_period=5 and time_window=Minutes"`
	Attribute      string              `json:"attribute" jsonschema:"required,description=The attribute to get values for"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the log attribute"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The filters to exclude from the log attribute"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the log attribute"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the log attribute"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get logs from"`
}

func GetLogAttributeValuesForIndividualAttributeHandler(arguments GetLogAttributeValuesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime := utils.CalculateTimeRange(arguments.TimeConfig)
	request := model.GetSingleLogSummaryRequest{
		LogSummaryRequest: model.LogSummaryRequest{
			StartTime:      startTime,
			EndTime:        endTime,
			Filters:        arguments.Filters,
			ExcludeFilters: arguments.ExcludeFilters,
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

	resp, err := utils.MakeMetoroAPIRequest("POST", "logsSummaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
