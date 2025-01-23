package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetMetricAttributeValuesForIndividualAttributeHandlerArgs struct {
	// Required: The metric name to get the summary for
	MetricName string `json:"metricName" jsonschema:"required,description=The name of the metric to get the attribute values for"`
	// The time period to get the attribute values for
	TimeConfig utils.TimeConfig `json:"timeConfig" jsonschema:"required,description=The time period to get the attribute values for. e.g. if you want to get the attribute values for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	// Environments is the environments to get the traces for. If empty, all environments will be included
	Environments []string `json:"environments" jsonschema:"description=The environments to get the attribute values from. If empty, all environments will be included"`
	// The attribute to get the summary for
	Attribute string `json:"attribute" jsonschema:"required,description=The attribute key to get the possible values for"`
}

type GetMetricSummaryForIndividualAttributeRequest struct {
	// Required: The metric name to get the summary for
	MetricName string `json:"metricName"`
	// Required: Start time of when to get the service summaries in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds since epoch
	EndTime int64 `json:"endTime"`
	// Environments is the environments to get the traces for. If empty, all environments will be included
	Environments []string `json:"environments"`
	// The attribute to get the summary for
	Attribute string `json:"attribute"`
}

type GetMetricSummaryForIndividualAttributeResponse struct {
	// The attribute values
	Attributes []string `json:"attribute"`
}

func GetMetricAttributeValuesForIndividualAttributeHandler(ctx context.Context, arguments GetMetricAttributeValuesForIndividualAttributeHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetMetricSummaryForIndividualAttributeRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		MetricName:   arguments.MetricName,
		Environments: arguments.Environments,
		Attribute:    arguments.Attribute,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "metricIndividualAttribute", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
