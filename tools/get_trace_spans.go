package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetTraceSpansHandlerArgs struct {
	TimeConfig   utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get trace spans for. e.g. if you want to get spans for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	TraceId      string           `json:"trace_id" jsonschema:"required,description=The traceId of the trace to get the associated spans. get_traces tool will return list of traceIds which should be used for this field."`
	Environments []string         `json:"environments" jsonschema:"description=The environments to get the spans for. If empty all environments will be included"`
}

type GetSpansForTraceRequest struct {
	// Required: Start time of when to get the traces in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the traces in seconds since epoch
	EndTime int64 `json:"endTime"`
	// Required: The traceId of the trace to get the associated spans.
	TraceId string `json:"traceId"`
	// The environments to get the traces for. If empty, all environments will be included
	Environments                   []string `json:"environments"`
	ShouldReturnNonMetoroEpbfSpans bool     `json:"shouldReturnNonMetoroEpbfSpans"`
}

func GetTraceSpansHandler(ctx context.Context, arguments GetTraceSpansHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetSpansForTraceRequest{
		StartTime:                      startTime,
		EndTime:                        endTime,
		TraceId:                        arguments.TraceId,
		Environments:                   arguments.Environments,
		ShouldReturnNonMetoroEpbfSpans: true,
	}

	body, err := getTraceSpansMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting trace spans: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getTraceSpansMetoroCall(ctx context.Context, request GetSpansForTraceRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling trace spans request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "spans", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
