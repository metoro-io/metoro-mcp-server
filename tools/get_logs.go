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
	Filters        map[string][]string `json:"filters" jsonschema:"description=Log attributes to restrict the search to. Keys are anded together and values in the keys are ORed.  e.g. {service.name: [/k8s/test/test /k8s/test/test2] namespace:[test]} will return all logs emited from (service.name = /k8s/test/test OR /k8s/test/test2) AND (namespace = test). Get the possible filter keys from the get_attribute_keys tool and possible values of a filter key from the get_attribute_values tool. If you are looking to get logs of a certain severity you should look up the log_level filter."`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=Log attributes to exclude from the search. Keys are anded together and values in the keys are ORed.  e.g. {service.name: [/k8s/test/test /k8s/test/test2] namespace:[test]} will return all logs emited from NOT ((service.name = /k8s/test/test OR /k8s/test/test2) AND (namespace = test)). Get the possible filter keys from the get_attribute_keys tool and possible values of a filter key from the get_attribute_values tool. If you are looking to get logs of a certain severity you should look up the log_level filter."`
	Regexes        []string            `json:"regexes" jsonschema:"description=Regexes to apply to the log messages. Only the logs with messages that match these regexes will be returned. Regexes are ANDed together. For example if you want to get logs with message that contains the word 'fish' and 'chips' you would set the regexes as ['fish' 'chips']. If you want to OR you should use the | operator in a single regex. regexes only match the body of the log so do not use this to match things like service names. If you want to get error logs use log_level filters instead."`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonshcema:"description=Regexes to exclude the log. Log messages that match these regexes will not be returned. Exclude regexes are ORed together. For example if you want to get logs with messages that do not contain the word 'fish' or 'chips' you would set the exclude regexes as ['fish' 'chips']. regexes only match the body of the log so do not use this to match things like service names. If you want to get error logs use log_level filters instead."`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get logs from. If empty logs from all environments will be returned"`
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
