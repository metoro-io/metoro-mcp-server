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

type GetKubernetesMetricHandlerArgs struct {
	// The time period to get the metrics for
	TimeConfig utils.TimeConfig `json:"timeConfig" jsonschema:"required,description=The time period to get the metrics for. e.g. if you want to get metrics for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	// The filters to apply to the kubernetes summary
	Filters map[string][]string `json:"filters" jsonschema:"description=The filters to apply to the kubernetes metrics. For example if you want to get metrics for a specific service you can pass in a filter like {'service.name': ['microservice_a']}"`
	// ExcludeFilters are filters that should be excluded
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to apply to the kubernetes metrics. For example if you want to get metrics for all services except microservice_a you can pass in {'service.name': ['microservice_a']}"`
	// Splits is a list of attributes to split the metrics by
	Splits []string `json:"splits" jsonschema:"description=List of attributes to split the metrics by. For example if you want to split the metrics by service you can pass in ['service.name']"`
	// The environments to get the kubernetes metrics for
	Environments []string `json:"environments" jsonschema:"description=The environments to get the kubernetes metrics from. If empty all environments will be included"`
	// Functions is the list of functions to apply to the metric
	Functions []model.MetricFunction `json:"functions" jsonschema:"description=The functions to apply to the metrics in order. Available functions are monotonicDifference valueDifference and customMathExpression"`
	// LimitResults is a flag to indicate if the results should be limited
	LimitResults bool `json:"limitResults" jsonschema:"description=Flag to indicate if the results should be limited"`
	// BucketSize is the size of each datapoint bucket in seconds
	BucketSize int64 `json:"bucketSize" jsonschema:"description=The size of each datapoint bucket in seconds"`
	// Aggregation is the operation to apply to the metrics
	Aggregation model.Aggregation `json:"aggregation" jsonschema:"required,description=The aggregation to apply to the metrics. Possible values are: count sum avg min max. If you want to get the number of pods use count."`
	// IsRate is a flag to indicate if the metric is a rate metric
	IsRate bool `json:"isRate" jsonschema:"description=Flag to indicate if the metric is a rate metric"`
	// JsonPath is a path to pull the json value from the metric
	JsonPath *string `json:"jsonPath" jsonschema:"description=Path to pull the json value from the metric"`
}

type GetKubernetesMetricRequest struct {
	// Required: Start time of when to get the service summaries in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds since epoch
	EndTime int64 `json:"endTime"`
	// The filters to apply to the kubernetes summary
	Filters map[string][]string `json:"filters"`
	// ExcludeFilters are filters that should be excluded
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// Splits is a list of attributes to split the metrics by
	Splits []string `json:"splits"`
	// The environments to get the kubernetes metrics for
	Environments []string `json:"environments"`
	// Functions is the list of functions to apply to the metric
	Functions []model.MetricFunction `json:"functions"`
	// LimitResults is a flag to indicate if the results should be limited
	LimitResults bool `json:"limitResults"`
	// BucketSize is the size of each datapoint bucket in seconds
	BucketSize int64 `json:"bucketSize"`
	// Aggregation is the operation to apply to the metrics
	Aggregation model.Aggregation `json:"aggregation"`
	// IsRate is a flag to indicate if the metric is a rate metric
	IsRate bool `json:"isRate"`
	// JsonPath is a path to pull the json value from the metric
	JsonPath *string `json:"jsonPath"`
}

func GetKubernetesMetricHandler(ctx context.Context, arguments GetKubernetesMetricHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetKubernetesMetricRequest{
		StartTime:      startTime,
		EndTime:        endTime,
		Filters:        arguments.Filters,
		ExcludeFilters: arguments.ExcludeFilters,
		Splits:         arguments.Splits,
		Environments:   arguments.Environments,
		Functions:      arguments.Functions,
		LimitResults:   arguments.LimitResults,
		BucketSize:     arguments.BucketSize,
		Aggregation:    arguments.Aggregation,
		IsRate:         arguments.IsRate,
		JsonPath:       arguments.JsonPath,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "kubernetesMetric", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
