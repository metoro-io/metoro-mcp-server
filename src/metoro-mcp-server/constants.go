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
