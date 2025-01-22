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

type GetMetricHandlerArgs struct {
	TimeConfig     utils.TimeConfig       `json:"time_config" jsonschema:"required,description=The time period to get the metric/timeseries data for. e.g. if you want to get the timeseries/metric data for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	MetricName     string                 `json:"metricName" jsonschema:"required,description=The name of the metric to use for getting the timeseries data"`
	Aggregation    model.Aggregation      `json:"aggregation" jsonschema:"required,description=The aggregation to apply to the metrics. Possible values are: sum / avg / min / max / count. The aggregation will be applied to every datapoint bucket. For example if the bucket size is 1 minute and the aggregation is sum then the sum of all datapoints in a minute will be returned"`
	Filters        map[string][]string    `json:"filters" jsonschema:"description=Filters to apply to the metrics. Only the metrics that match these filters will be returned. Get the possible filter keys and values the get_metric_attributes tool. e.g. {service_name: [/k8s/namespaceX/serviceX]} should return metrics for serviceX in namespaceX"`
	ExcludeFilters map[string][]string    `json:"excludeFilters" jsonschema:"description=Filters to exclude the metrics. Metrics matching the exclude filters will not be returned. Get the possible exclude filter keys and values from the get_metric_attributes tool. e.g. {service_name: [/k8s/namespaceX/serviceX]} should exclude metrics from serviceX in namespaceX"`
	Splits         []string               `json:"splits" jsonschema:"description=Splits will allow you to group/split metrics by an attribute. This is useful if you would like to see the breakdown of a particular metric by an attribute. For example if you want to see the breakdown of the metric by service_name you would set the splits as ['service_name']"`
	Functions      []model.MetricFunction `json:"functions" jsonschema:"description=The functions to apply to the metric. Available functions are monotonicDifference which will calculate the difference between the current and previous value of the metric (negative values will be set to 0) and valueDifference which will calculate the difference between the current and previous value of the metric or MathExpression e.g. a / 60"`
	BucketSize     int64                  `json:"bucketSize" jsonschema:"description=The size of each datapoint bucket in seconds if not provided metoro will select the best bucket size for the given duration for performance and clarity"`
}

func GetMetricHandler(ctx context.Context, arguments GetMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetMetricRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		MetricName:     arguments.MetricName,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Splits:         arguments.Splits,
		Aggregation:    arguments.Aggregation,
		Functions:      arguments.Functions,
		BucketSize:     arguments.BucketSize,
	}

	body, err := getMetricMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting metric: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getMetricMetoroCall(ctx context.Context, request model.GetMetricRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling metric request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "metric", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
