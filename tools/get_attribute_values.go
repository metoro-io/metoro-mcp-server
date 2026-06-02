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
	Type       model.MetricType `json:"type" jsonschema:"required,description=The type of telemetry data to get the attribute keys and values for. Either 'logs' or 'trace' or 'metric' or 'kubernetes_resource'. The alias 'kubernetes_resources' is also accepted."`
	TimeConfig utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to use while getting the possible values of log attributes. e.g. if you want to get values for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Attribute  string           `json:"attribute" jsonschema:"required,description=The attribute key to get the possible values for. Possible values for attribute should be obtained from get_attribute_keys tool call for the same type"`
	MetricName string           `json:"metricName" jsonschema:"description=REQUIRED IF THE TYPE IS 'metric'. The name of the metric to get the possible attribute keys and values."`
	Filters    []model.Filter   `json:"filters" jsonschema:"description=The filters to apply before getting the possible values. For example if you want to get the possible values for an attribute key where the environment is X you would set the Filters as [{key: 'environment' values: ['X']}]"`
}

func GetAttributeValuesHandler(ctx context.Context, arguments GetAttributeValuesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	// Convert Filter slice to map format for internal API
	metricType := normalizeMetricType(arguments.Type)
	filters := model.FiltersToMap(arguments.Filters)
	attribute := arguments.Attribute
	if metricType == model.KubernetesResource {
		filters = normalizeKubernetesAttributeMap(filters)
		attribute = normalizeKubernetesAttribute(attribute)
	}

	request := model.GetAttributeValuesRequest{
		Type:      metricType,
		Attribute: attribute,
	}

	switch metricType {
	case model.Logs:
		err = CheckAttributes(ctx, metricType, filters, map[string][]string{}, []string{}, nil)
		if err != nil {
			return nil, err
		}
		modelRequest := model.LogSummaryRequest{
			StartTime: startTime,
			EndTime:   endTime,
			Filters:   filters,
		}
		request.Logs = &modelRequest
		break
	case model.Trace:
		err = CheckAttributes(ctx, metricType, filters, map[string][]string{}, []string{}, nil)
		if err != nil {
			return nil, err
		}
		modelRequest := model.TracesSummaryRequest{
			StartTime: startTime,
			EndTime:   endTime,
			Filters:   filters,
		}
		request.Trace = &modelRequest
		break
	case model.Metric:
		err = CheckAttributes(ctx, metricType, filters, map[string][]string{}, []string{}, &model.GetMetricAttributesRequest{
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
			Environments: filters["environment"],
		}
		request.Metric = &modelRequest
		break
	case model.KubernetesResource:
		err = CheckAttributes(ctx, metricType, filters, map[string][]string{}, []string{attribute}, nil)
		if err != nil {
			return nil, err
		}
		modelRequest := model.GetKubernetesResourceRequest{
			StartTime: startTime,
			EndTime:   endTime,
			Filters:   filters,
		}
		request.Kubernetes = &modelRequest
		break
	default:
		return nil, fmt.Errorf("invalid type: %v", metricType)
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
