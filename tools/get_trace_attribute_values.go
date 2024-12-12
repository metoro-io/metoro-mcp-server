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

type GetTraceAttributeValuesHandlerArgs struct {
	TimeConfig   utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get trace attributes for. e.g. if you want to get attributes for the last 5 minutes, you would set time_period=5 and time_window=Minutes"`
	Attribute    string           `json:"attribute" jsonschema:"required, description=The name of the attribute to get values for"`
	ServiceNames []string         `json:"serviceNames" jsonschema:"description=The service names to get attribute values for"`
	//  TODO: I don't think we need these two fields for the LLM tool
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the traces. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the traces. Traces matching the exclude traces will not be returned. it is a map of filter keys to array values where array values are ORed when the filters are applied.e.g. key for service name is service.name"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the trace's endpoints. Traces with endpoints matching regexes will be returned"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The regexes to exclude from the trace's endpoints. Traces with endpoints matching regexes will be excluded"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get traces from. If empty, traces from all environments will be returned"`
}

func GetTraceAttributeValuesForIndividualAttributeHandler(arguments GetTraceAttributeValuesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	hourAgo := now.Add(-1 * time.Hour)
	request := model.GetSingleTraceSummaryRequest{
		TracesSummaryRequest: model.TracesSummaryRequest{
			StartTime:      hourAgo.Unix(),
			EndTime:        now.Unix(),
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

	resp, err := utils.MakeMetoroAPIRequest("POST", "tracesSummaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
