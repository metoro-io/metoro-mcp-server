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

type GetMetricAttributesHandlerArgs struct {
	TimeConfig       utils.TimeConfig `json:"timeConfig" jsonschema:"required,description=The time period to get the possible values of metric attributes for. e.g. if you want to get the possible values for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	MetricName       string           `json:"metricName" jsonschema:"required,description=The name of the metric to get the possible attribute keys and values."`
	FilterAttributes []model.Filter   `json:"filterAttributes" jsonschema:"description=The attributes to filter the metric attributes by before getting the possible values. For example if you want to get the possible keys and values where the environment is X you would set the filterAttributes as [{key: 'environment' values: ['X']}]"`
}

func GetMetricAttributesHandler(ctx context.Context, arguments GetMetricAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	// Convert Filter slice to map format for internal API
	filterAttributes := model.FiltersToMap(arguments.FilterAttributes)

	request := model.MetricAttributesRequest{
		StartTime:        startTime,
		EndTime:          endTime,
		MetricName:       arguments.MetricName,
		FilterAttributes: filterAttributes,
	}
	response, err := getMetricAttributesMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}

func getMetricAttributesMetoroCall(ctx context.Context, request model.MetricAttributesRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return utils.MakeMetoroAPIRequest("POST", "metricAttributes", bytes.NewBuffer(jsonData), utils.GetAPIRequirementsFromRequest(ctx))
}
