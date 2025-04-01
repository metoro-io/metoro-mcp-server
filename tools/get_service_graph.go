package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetServiceGraphHandlerArgs struct {
	TimeConfig   utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get the service graph for. e.g. if you want to get the graph for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	ServiceName  string           `json:"serviceName" jsonschema:"required,description=The name of the service to get the graph for"`
	Environments []string         `json:"environments" jsonschema:"description=The environments to get the service graph for. If empty all environments will be used."`
}

type GetServiceGraphRequest struct {
	// If Environments is not empty, only services that are in the list will be included in the graph.
	// If Environments is empty, all services will be included in the graph.
	Environments []string `json:"environments"`
	// If InitialServices is not empty, only services that are in the list will be included in the graph.
	// If InitialServices is empty, all services will be included in the graph.
	InitialServices []string `json:"initialServices"`
	// If EndingServices is not empty, only services that are in the list will be included in the graph.
	// If EndingServices is empty, all services will be included in the graph.
	EndingServices []string `json:"endingServices"`
	// StartTime is the start time of the graph in seconds since epoch
	StartTime int64 `json:"startTime"`
	// EndTime is the end time of the graph in seconds since epoch
	EndTime int64 `json:"endTime"`
	// The filters to apply to the traces, so for example, if you want to get traces for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`
	// ExcludeFilters are filters that should be excluded from the traces
	// For example, if you want to get traces for all services except microservice_a you can pass in
	// {"service_name": ["microservice_a"]}
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// Regexes are used to filter traces based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter traces based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`
}

func GetServiceGraphHandler(ctx context.Context, arguments GetServiceGraphHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetServiceGraphRequest{
		StartTime:       startTime,
		EndTime:         endTime,
		Environments:    arguments.Environments,
		InitialServices: []string{arguments.ServiceName},
		EndingServices:  []string{arguments.ServiceName},
	}

	body, err := getServiceGraphMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting service graph: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getServiceGraphMetoroCall(ctx context.Context, request GetServiceGraphRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling service graph request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "serviceGraph", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
