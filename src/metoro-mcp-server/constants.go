package main

const METORO_API_URL_ENV_VAR = "METORO_API_URL"
const METORO_AUTH_TOKEN_ENV_VAR = "METORO_AUTH_TOKEN"

// TODO: This file should be replaced if we can import the types from Metoro repo directly.
// These are just duplicates at the moment. If updated in Metoro repository, it should also be updated here!

type GetLogsRequest struct {
	// Required: Start time of when to get the logs in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the logs in seconds since epoch
	EndTime int64 `json:"endTime"`
	// The filters to apply to the logs, so for example, if you want to get logs for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`
	// ExcludeFilters are filters that should be excluded from the logs
	// For example, if you want to get logs for all services except microservice_a you can pass in
	// {"service_name": ["microservice_a"]}
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// Previous page endTime in nanoseconds, used to get the next page of logs if there are more logs than the page size
	// If omitted, the first page of logs will be returned
	PrevEndTime *int64 `json:"prevEndTime"`
	// Regexes are used to filter logs based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter logs based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`
	Ascending      bool     `json:"ascending"`
	// The cluster/environments to get the logs for. If empty, all clusters will be included
	Environments []string `json:"environments"`
}

type GetTracesRequest struct {
	ServiceNames   []string            `json:"serviceNames"`
	StartTime      int64               `json:"startTime"`
	EndTime        int64               `json:"endTime"`
	Filters        map[string][]string `json:"filters"`
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	PrevEndTime    *int64              `json:"prevEndTime"`
	Regexes        []string            `json:"regexes"`
	ExcludeRegexes []string            `json:"excludeRegexes"`
	Ascending      bool                `json:"ascending"`
	Environments   []string            `json:"environments"`
}

type Aggregation string

const (
	AggregationSum   Aggregation = "sum"
	AggregationAvg   Aggregation = "avg"
	AggregationMax   Aggregation = "max"
	AggregationMin   Aggregation = "min"
	AggregationCount Aggregation = "count"
	AggregationP50   Aggregation = "p50"
	AggregationP90   Aggregation = "p90"
	AggregationP95   Aggregation = "p95"
	AggregationP99   Aggregation = "p99"

	// Only for trace metrics
	AggregationRequestSize  Aggregation = "requestSize"
	AggregationResponseSize Aggregation = "responseSize"
	AggregationTotalSize    Aggregation = "totalSize"
)

type MetricFunction struct {
	ID string `json:"id"`
	// The type of the function
	FunctionType FunctionType `json:"functionType"`
	// The payload of the function
	// TODO: If we have more payloads this can be an interface but for now its a math expression since its the only payload.
	FunctionPayload MathExpression `json:"functionPayload"`
}

type MathExpression struct {
	Variables  []string `json:"variables"`
	Expression string   `json:"expression"`
}

type FunctionType string

const (
	MonotonicDifference  FunctionType = "monotonicDifference"
	ValueDifference      FunctionType = "valueDifference"
	CustomMathExpression FunctionType = "customMathExpression"
)

type GetMetricRequest struct {
	// MetricName is the name of the metric to get
	MetricName string `json:"metricName"`
	// Required: Start time of when to get the logs in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the logs in seconds since epoch
	EndTime int64 `json:"endTime"`
	// The filters to apply to the logs, so for example, if you want to get logs for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`
	// The filters to exclude from the logs, so for example, if you want to exclude logs for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// Splits is a list of attributes to split the metrics by, for example, if you want to split the metrics by service
	// you can pass in a list like ["service_name"]
	Splits []string `json:"splits"`
	// Aggregation is the operation to apply to the metrics, for example, if you want to sum the metrics you can pass in "sum"
	Aggregation Aggregation `json:"aggregation"`
	// IsRate is a flag to indicate if the metric is a rate metric
	IsRate bool `json:"isRate"`
	// Functions is the list of functions to apply to the metric, in the same order that they appear in this array!!
	Functions []MetricFunction `json:"functions"`
	// LimitResults is a flag to indicate if the results should be limited.
	LimitResults bool `json:"limitResults"`
	// BucketSize is the size of each datapoint bucket in seconds
	BucketSize int64 `json:"bucketSize"`
}

type GetProfileRequest struct {
	// Required: ServiceName to get profiling for
	ServiceName string `json:"serviceName"`

	// Optional: ContainerNames to get profiling for
	ContainerNames []string `json:"containerNames"`

	// Required: Timestamp to get profiling after this time
	// Seconds since epoch
	StartTime int64 `json:"startTime"`

	// Required: Timestamp to get profiling this time
	// Seconds since epoch
	EndTime int64 `json:"endTime"`
}
type GetTraceMetricRequest struct {
	// Required: Start time of when to get the logs in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the logs in seconds since epoch
	EndTime int64 `json:"endTime"`

	// Optional: The name of the service to get the trace metrics for
	// Acts as an additional filter
	ServiceNames []string `json:"serviceNames"`

	// The filters to apply to the logs, so for example, if you want to get logs for a specific service
	//you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`

	// The exclude filters to apply to the logs, so for example, if you want to exclude logs for a specific service
	//you can pass in a filter like {"service_name": ["microservice_a"]}
	ExcludeFilters map[string][]string `json:"excludeFilters"`

	// Regexes are used to filter traces based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter traces based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`

	// Splts is a list of attributes to split the metrics by, for example, if you want to split the metrics by service
	// you can pass in a list like ["service_name"]
	Splits []string `json:"splits"`

	// Functions is the array of function to apply to the trace metrics,
	//for example, if you want to get the monotonic difference between count of traces each minute.
	// Functions are applied in the same order that they appear in this array
	Functions []MetricFunction `json:"functions"`

	// Aggregate to apply to trace metrics, for example, if you want to sum the metrics you can pass in "sum"
	Aggregate Aggregation `json:"aggregate"`

	// Environments is a list of environments to filter the traces by. If empty, all environments will be included
	Environments []string `json:"environments"`

	// LimitResults is a flag to indicate if the results should be limited.
	LimitResults bool `json:"limitResults"`

	// BucketSize is the size of each datapoint bucket in seconds
	BucketSize int64 `json:"bucketSize"`
}

type GetK8sEventsRequest struct {
	// Required: Start time of when to get the k8s events in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the k8s events in seconds since epoch
	EndTime int64 `json:"endTime"`
	// The filters to apply to the k8s events, so for example, if you want to get k8s events for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`
	// ExcludeFilters are filters that should be excluded from the k8s events
	// For example, if you want to get k8s events for all services except microservice_a you can pass in
	// {"service_name": ["microservice_a"]}
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// Previous page endTime in nanoseconds, used to get the next page of k8s events if there are more k8s events than the page size
	// If omitted, the first page of k8s events will be returned
	PrevEndTime *int64 `json:"prevEndTime"`
	// Regexes are used to filter k8s events based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter k8s events based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`
	// Ascending is a flag to determine if the k8s events should be returned in ascending order
	Ascending bool `json:"ascending"`
	// Environments is the environments to get the k8s events for
	Environments []string `json:"environments"`
}

type GetSingleK8sEventSummaryRequest struct {
	GetK8sEventsRequest
	// The attribute to get the summary for
	Attribute string `json:"attribute"`
}

type GetK8sEventMetricsRequest struct {
	// Required: Start time of when to get the logs in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the logs in seconds since epoch
	EndTime int64 `json:"endTime"`

	// The filters to apply to the logs, so for example, if you want to get logs for a specific service
	//you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`

	// The exclude filters to apply to the logs, so for example, if you want to exclude logs for a specific service
	//you can pass in a filter like {"service_name": ["microservice_a"]}
	ExcludeFilters map[string][]string `json:"excludeFilters"`

	// Regexes are used to filter k8s events based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter k8s events based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`

	// Splts is a list of attributes to split the metrics by, for example, if you want to split the metrics by service
	// you can pass in a list like ["service_name"]
	Splits []string `json:"splits"`

	// OnlyNumRequests is a flag to only get the number of requests, this is a much faster query
	OnlyNumRequests bool `json:"onlyNumRequests"`

	// Environments is a list of environments to filter the k8s events by. If empty, all environments will be included
	Environments []string `json:"environments"`
}

type GetPodsRequest struct {
	// Required: Timestamp to get metadata updates after this time
	StartTime int64 `json:"startTime"`

	// Required: Timestamp to get metadata updates before this time
	EndTime int64 `json:"endTime"`

	// Optional: Environment to filter the pods by. If not provided, all environments are considered
	Environments []string `json:"environments"`

	// Optional: ServiceName to get metadata updates. One of ServiceName or NodeName is required
	ServiceName string `json:"serviceName"`

	// Optional: NodeName to get metadata updates. One of ServiceName or NodeName is required
	NodeName string `json:"nodeName"`
}

type LogSummaryRequest struct {
	// Required: Start time of when to get the service summaries in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds since epoch
	EndTime int64 `json:"endTime"`
	// The filters to apply to the log summary, so for example, if you want to get logs for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`

	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// RegexFilter is a regex to filter the logs
	Regexes        []string `json:"regexes"`
	ExcludeRegexes []string `json:"excludeRegexes"`
	// The cluster/environments to get the logs for. If empty, all clusters will be included
	Environments []string `json:"environments"`
}

type GetSingleLogSummaryRequest struct {
	LogSummaryRequest
	// The attribute to get the summary for
	Attribute string `json:"attribute"`
}

type GetAllNodesRequest struct {
	// StartTime Required: Start time of when to get the nodes in seconds since epoch
	StartTime int64 `json:"startTime"`
	// EndTime Required: End time of when to get the nodes in seconds since epoch
	EndTime int64 `json:"endTime"`
	// Environments The cluster/environments to get the nodes for. If empty, all clusters will be included
	Environments []string `json:"environments"`
	// Filters The filters to apply to the nodes, so for example, if you want to get subset of nodes that have a specific label
	Filters map[string][]string `json:"filters"`
	// ExcludeFilters are filters that should be excluded from the nodes
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	// Splits is a list of attributes to split the nodes by, for example, if you want to split the nodes a label
	Splits []string `json:"splits"`
}
