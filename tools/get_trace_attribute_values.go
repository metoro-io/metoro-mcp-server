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

type GetTraceAttributeValuesHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to use for getting the possible values for a trace attribute key. e.g. if you want to get possible trace attribute values for key x for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Attribute      string              `json:"attribute" jsonschema:"required,description=The name of the attribute key to get the possible values for"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply before getting the possible values. For example if you want to get the possible values for attribute key service.name where the environment is X you would set the Filters as {environment: [X]}"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate possible values an attribute can take. Traces matching the exclude filters will not be returned. For example if you want the possible values for attribute key service.name where the attribute environment is not X then you would set the ExcludeFilters as {environment: [X]}"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the trace endpoint. Only the attribute values (for a given attribute key) of trace endpoint that match these regexes will be returned. For example if you want the possible values for attribute key service.name where the trace endpoint contains the word 'get' you would set the regexes as ['get']"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The exclude regexes to apply to the trace endpoint. The attribute values (for a given attribute key) of trace endpoint that match these regexes will not be returned. For example if you want the possible values for attribute key service.name where the trace endpoint does not contain the word 'get' you would set the exclude regexes as ['get']"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get traces from. If empty traces from all environments will be returned"`
}

func GetTraceAttributeValuesForIndividualAttributeHandler(ctx context.Context, arguments GetTraceAttributeValuesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetSingleTraceSummaryRequest{
		TracesSummaryRequest: model.TracesSummaryRequest{
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

	resp, err := utils.MakeMetoroAPIRequest("POST", "tracesSummaryIndividualAttribute", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
