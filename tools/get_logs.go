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
	Filters        map[string][]string `json:"attributeFilters" jsonschema:"description=You must use get_attribute_keys and get_attribute_values before setting this. Log attributes to restrict the search to. Keys are anded together and values in the keys are ORed.  e.g. {service.name: [/k8s/test/test /k8s/test/test2] namespace:[test]} will return all logs emited from (service.name = /k8s/test/test OR /k8s/test/test2) AND (namespace = test). Get the possible filter keys from the get_attribute_keys tool and possible values of a filter key from the get_attribute_values tool. If you are looking to get logs of a certain severity you should look up the log_level filter."`
	ExcludeFilters map[string][]string `json:"attributeExcludeFilters" jsonschema:"description=You must use get_attribute_keys and get_attribute_values before setting this.Log attributes to exclude from the search. Keys are anded together and values in the keys are ORed.  e.g. {service.name: [/k8s/test/test /k8s/test/test2] namespace:[test]} will return all logs emited from NOT ((service.name = /k8s/test/test OR /k8s/test/test2) AND (namespace = test)). Get the possible filter keys from the get_attribute_keys tool and possible values of a filter key from the get_attribute_values tool. If you are looking to get logs of a certain severity you should look up the log_level filter."`
	Regex          string              `json:"regex" jsonschema:"description=Regex to apply to the log search re2 format. Any match in the log message will cause it to be returned. Use the filters parameter log_level if you want to look for logs of a certain severity"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get logs from. If empty logs from all environments will be returned"`
}

func GetLogsHandler(ctx context.Context, arguments GetLogsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	var regexes = []string{}
	if arguments.Regex != "" {
		regexes = append(regexes, arguments.Regex)
	}

	err = CheckAttributes(ctx, model.Logs, arguments.Filters, arguments.ExcludeFilters, []string{}, nil)
	if err != nil {
		return nil, err
	}
	limit := 20

	request := model.GetLogsRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Regexes:        regexes,
		Environments:   arguments.Environments,
		ExportLimit:    &limit,
	}

	resp, err := getLogsMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting logs: %v", err)
	}
	respTrimmed, err := trimLogsResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error trimming logs: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(respTrimmed)))), nil
}

func getLogsMetoroCall(ctx context.Context, request model.GetLogsRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling logs request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "logs", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}

func trimLogsResponse(response []byte) ([]byte, error) {
	var logsResponse model.GetLogsResponse
	err := json.Unmarshal(response, &logsResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling logs response: %v", err)
	}

	// Trim every log entry to only include the first 1000 characters of the message
	// This is to prevent excessively long log messages from blowing up the context.
	logLineLengthLimit := 1000
	for i := range logsResponse.Logs {
		if len(logsResponse.Logs[i].Message) > logLineLengthLimit {
			logsResponse.Logs[i].Message = logsResponse.Logs[i].Message[:logLineLengthLimit]
		}
	}

	return json.Marshal(logsResponse)
}
