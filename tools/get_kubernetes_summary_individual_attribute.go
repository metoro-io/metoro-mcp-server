package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetKubernetesSummaryForIndividualAttributeHandlerArgs struct {
	// The attribute to get
	Attribute string `json:"attribute" jsonschema:"required,description=The attribute key to get the summary for"`
	// The time period to get the summary for
	TimeConfig utils.TimeConfig `json:"timeConfig" jsonschema:"required,description=The time period to get the summary for. e.g. if you want to get the summary for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	// The filters to apply to the summary
	Filters map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the kubernetes summary. For example, if you want to get summary for a specific service you can pass in a filter like {'service.name': ['microservice_a']}"`
	// ExcludeFilters are filters that should be excluded
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to apply to the kubernetes summary. For example, if you want to get summary for all services except microservice_a you can pass in {'service.name': ['microservice_a']}"`
	// Regexes are used to filter based on a regex inclusively
	Regexes []string `json:"regexes" jsonschema:"description=Regexes to filter the kubernetes summary inclusively"`
	// ExcludeRegexes are used to filter based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes" jsonschema:"description=Regexes to filter the kubernetes summary exclusively"`
	// The environments to get the summary for
	Environments []string `json:"environments" jsonschema:"description=The environments to get the kubernetes summary from. If empty, all environments will be included"`
}

type GetKubernetesSummaryForIndividualAttributeRequest struct {
	// The attribute to get
	Attribute string `json:"attribute"`
	// Required: Start time of when to get the logs in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the logs in seconds since epoch
	EndTime int64 `json:"endTime"`
	// The filters to apply to the logs
	Filters map[string][]string `json:"filters"`
	// The exclude filters to apply to the logs
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// Regexes are used to filter based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`
	// The environments to get the summary for
	Environments []string `json:"environments"`
}

func GetKubernetesSummaryForIndividualAttributeHandler(ctx context.Context, arguments GetKubernetesSummaryForIndividualAttributeHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetKubernetesSummaryForIndividualAttributeRequest{
		Attribute:      arguments.Attribute,
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

	resp, err := utils.MakeMetoroAPIRequest("POST", "kubernetesSummaryIndividualAttribute", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
