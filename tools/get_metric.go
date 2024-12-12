package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetMetricHandlerArgs struct {
	TimeConfig     utils.TimeConfig       `json:"time_config" jsonschema:"required,description=The time period to get the metric/timeseries data for. e.g. if you want to get the timeseries/metric data for the last 5 minutes, you would set time_period=5 and time_window=Minutes"`
	MetricName     string                 `json:"metricName" jsonschema:"required,description=The name of the metric to get"`
	Aggregation    model.Aggregation      `json:"aggregation" jsonschema:"required,description=The aggregation to apply to the metric. e.g. sum, avg, min, max, count"`
	Filters        map[string][]string    `json:"filters" jsonschema:"description=Filters to apply to the metric. Metrics matching the filters will be returned"`
	ExcludeFilters map[string][]string    `json:"excludeFilters" jsonschema:"description=Exclude filters to apply to the metric. Metrics matching the exclude filters will not be returned"`
	Splits         []string               `json:"splits" jsonschema:"description=The splits to apply to the metric. Metrics will be split by the given keys"`
	Functions      []model.MetricFunction `json:"functions" jsonschema:"description=The functions to apply to the metric"`
	LimitResults   bool                   `json:"limitResults" jsonschema:"description=If true, the results will be limited to improve performance"`
	BucketSize     int64                  `json:"bucketSize" jsonschema:"description=The size of each datapoint bucket in seconds, if not provided metoro will select the best bucket size for the given duration for performance and clarity"`
}

func GetMetricHandler(arguments GetMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
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
		LimitResults:   arguments.LimitResults,
		BucketSize:     arguments.BucketSize,
	}

	body, err := getMetricMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting metric: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getMetricMetoroCall(request model.GetMetricRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling metric request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "metric", bytes.NewBuffer(requestBody))
}
