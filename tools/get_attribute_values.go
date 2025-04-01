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

type GetAttributeValuesHandlerArgs struct {
	Type           model.MetricType    `json:"type" jsonschema:"required,description=The type of telemetry data to get the attribute keys and values for. Either 'logs' or 'trace' or 'metric'."`
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to use while getting the possible values of log attributes. e.g. if you want to get values for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Attribute      string              `json:"attribute" jsonschema:"required,description=The attribute key to get the possible values for. Possible values for attribute should be obtained from get_attribute_keys tool call for the same type"`
	MetricName     string              `json:"metricName" jsonschema:"description=REQUIRED IF THE TYPE IS 'metric'. The name of the metric to get the possible attribute keys and values."`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply before getting the possible values. For example if you want to get the possible values for attribute key service.name where the environment is X you would set the Filters as {environment: [X]}"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate possible values an attribute can take. Log attributes matching the exclude filters will not be returned. For example if you want the possible values for attribute key service.name where the attribute environment is not X then you would set the ExcludeFilters as {environment: [X]}"`
}

func GetAttributeValuesHandler(ctx context.Context, arguments GetAttributeValuesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	err = CheckAttributes(ctx, arguments.Type, arguments.Filters, arguments.ExcludeFilters, []string{}, nil)
	if err != nil {
		return nil, err
	}

	request := model.GetAttributeValuesRequest{
		Type:      arguments.Type,
		Attribute: arguments.Attribute,
	}

	switch arguments.Type {
	case model.Logs:
		err = CheckAttributes(ctx, arguments.Type, arguments.Filters, arguments.ExcludeFilters, []string{}, nil)
		if err != nil {
			return nil, err
		}
		modelRequest := model.LogSummaryRequest{
			StartTime:      startTime,
			EndTime:        endTime,
			Filters:        arguments.Filters,
			ExcludeFilters: arguments.ExcludeFilters,
		}
		request.Logs = &modelRequest
		break
	case model.Trace:
		err = CheckAttributes(ctx, arguments.Type, arguments.Filters, arguments.ExcludeFilters, []string{}, nil)
		if err != nil {
			return nil, err
		}
		modelRequest := model.TracesSummaryRequest{
			StartTime:      startTime,
			EndTime:        endTime,
			Filters:        arguments.Filters,
			ExcludeFilters: arguments.ExcludeFilters,
		}
		request.Trace = &modelRequest
		break
	case model.Metric:
		err = CheckAttributes(ctx, arguments.Type, arguments.Filters, arguments.ExcludeFilters, []string{}, &model.GetMetricAttributesRequest{
			StartTime:  startTime,
			EndTime:    endTime,
			MetricName: arguments.MetricName,
		})
		if err != nil {
			return nil, err
		}
		modelRequest := model.GetMetricAttributesRequest{
			StartTime:    startTime,
			EndTime:      endTime,
			MetricName:   arguments.MetricName,
			Environments: arguments.Filters["environment"],
		}
		request.Metric = &modelRequest
		break
	//case model.KubernetesResource:
	//
	//	modelRequest := model.GetKubernetesResourceRequest{
	//		StartTime:      startTime,
	//		EndTime:        endTime,
	//		Filters:        arguments.Filters,
	//		ExcludeFilters: arguments.Filters,
	//	}
	//	request.Kubernetes = &modelRequest
	//	break
	default:
		return nil, fmt.Errorf("invalid type: %v", arguments.Type)
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "metrics/attribute/values", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))

	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
