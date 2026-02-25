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

type GetMetricNamesHandlerArgs struct {
	TimeConfig       utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get metric names for. e.g. if you want to get metric names from the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	FuzzyStringMatch string           `json:"fuzzy_string_match,omitempty" jsonschema:"description=Optional fuzzy match string to search metric names"`
	Environments     []string         `json:"environments" jsonschema:"description=Environments to get metrics names from. If empty all environments will be used."`
}

func GetMetricNamesHandler(ctx context.Context, arguments GetMetricNamesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := model.FuzzyMetricsRequest{
		StartTime:        startTime,
		EndTime:          endTime,
		MetricFuzzyMatch: arguments.FuzzyStringMatch,
		Environments:     arguments.Environments,
	}
	response, err := getMetricNamesMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}

func getMetricNamesMetoroCall(ctx context.Context, request model.FuzzyMetricsRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return utils.MakeMetoroAPIRequest("POST", "fuzzyMetricsNames", bytes.NewBuffer(jsonData), utils.GetAPIRequirementsFromRequest(ctx))
}
