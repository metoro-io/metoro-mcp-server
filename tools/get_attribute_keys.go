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

type GetAttributeKeysHandlerArgs struct {
	Type       model.MetricType `json:"type" jsonschema:"required,description=The type of attribute keys to get. Either 'logs' or 'trace' or 'metric' or 'kubernetes_resource'"`
	TimeConfig utils.TimeConfig `json:"timeConfig" jsonschema:"required,description=The time period to get the possible attribute keys for a metric. This is required if type is 'metric' e.g. if you want to get the possible values for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	MetricName string           `json:"metricName" jsonschema:"description=The name of the metric to get the possible attribute keys for. This is required if type is 'metric'"`
}

func GetAttributeKeysHandler(ctx context.Context, arguments GetAttributeKeysHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	metricAttr := model.GetMetricAttributesRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		MetricName:   arguments.MetricName,
		Environments: []string{}, // TODO: Add environments to the request if needed. For now, we are not using it as I don't think its needed.
	}

	request := model.MultiMetricAttributeKeysRequest{
		Type:   string(arguments.Type),
		Metric: &metricAttr,
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "metrics/attributes", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))

	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
