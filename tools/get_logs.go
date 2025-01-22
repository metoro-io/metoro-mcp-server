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

type GetLogsHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get the logs for. e.g. if you want the get the logs for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=Filters to apply to the logs. Only the logs that match these filters will be returned. Get the possible filter keys from the get_log_attributes tool and possible values of a filter key from the get_log_attribute_values_for_individual_attribute tool. e.g. {service.name: [/k8s/namespaceX/serviceX]} should return logs for serviceX in namespaceX"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=Filters to exclude the logs. Logs matching the exclude filters will not be returned. Get the possible exclude filter keys from the get_log_attributes tool and possible exclude filter values from the get_log_attribute_values_for_individual_attribute tool (for a key). e.g. {service.name: [/k8s/namespaceX/serviceX]} should exclude logs for serviceX in namespaceX"`
	Regexes        []string            `json:"regexes" jsonschema:"description=Regexes to apply to the log messages. Only the logs with messages that match these regexes will be returned. Regexes are ORed together. For example if you want to get logs with message that contains the word 'error' or 'warning' you would set the regexes as ['error' 'warning']"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonshcema:"description=Regexes to exclude the log. Log messages that match these regexes will not be returned. Exclude regexes are AND together. For example if you want to get logs with messages that do not contain the word 'error' or 'warning' you would set the exclude regexes as ['error' 'warning']"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get logs from. If empty, logs from all environments will be returned"`
}

func GetLogsHandler(ctx context.Context, arguments GetLogsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := model.GetLogsRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        arguments.Regexes,
		ExcludeRegexes: arguments.ExcludeRegexes,
		Environments:   arguments.Environments,
	}

	resp, err := getLogsMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting logs: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func getLogsMetoroCall(ctx context.Context, request model.GetLogsRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling logs request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "logs", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
