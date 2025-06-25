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

type GetTracesHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get traces for. e.g. if you want to get traces for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time. Try to use a time period 1 hour or less unless its requested."`
	Filters        map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the traces. Only the traces that match these filters will be returned. You have to get the possible filter keys from the get_attribute_keys tool and possible values of a filter key from the get_attribute_values tool. DO NOT GUESS THE FILTER KEYS OR VALUES. Multiple filter keys are ANDed together and values for a filter key are ORed together"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate the traces. Traces matching the exclude traces will not be returned. You have to get the possible exclude filter keys from the get_attribute_keys tool and possible value for the key from the get_attribute_values tool. DO NOT GUESS THE FILTER KEYS OR VALUES. Multiple keys are ORed together and values for a filter key are ANDed together"`
}

func GetTracesHandler(ctx context.Context, arguments GetTracesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	err = CheckAttributes(ctx, model.Trace, arguments.Filters, arguments.ExcludeFilters, []string{}, nil)
	if err != nil {
		return nil, err
	}

	limit := 20

	request := model.GetTracesRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Limit:          &limit,
	}

	body, err := getTracesMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting traces: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getTracesMetoroCall(ctx context.Context, request model.GetTracesRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling traces request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "traces", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
