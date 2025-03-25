package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"

	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetMultiMetricHandlerArgs struct {
	TimeConfig utils.TimeConfig            `json:"time_config" jsonschema:"required,description=The time period to get the metric/timeseries data for. e.g. if you want to get the timeseries/metric data for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Metrics    []model.SingleMetricRequest `json:"metrics" jsonschema:"required,description=Array of metrics to get the timeseries data for"`
	Formulas   []model.Formula             `json:"formulas" jsonschema:"description=Optional formulas to combine metrics. Formula should only consist of formulaIdentifier of the metrics in  the metrics array. e.g. a + b + c if a b c appears in the formulaIdentifier of the metrics array. You can do * / + - && || ?: operations"`
}

func GetMultiMetricHandler(ctx context.Context, arguments GetMultiMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetMultiMetricRequest{
		StartTime: startTime,
		EndTime:   endTime,
		Metrics:   arguments.Metrics,
		Formulas:  arguments.Formulas,
	}

	body, err := getMultiMetricMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting metric: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getMultiMetricMetoroCall(ctx context.Context, request model.GetMultiMetricRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling metric request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "metrics", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
