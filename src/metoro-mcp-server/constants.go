package main

const METORO_API_URL_ENV_VAR = "METORO_API_URL"
const METORO_AUTH_TOKEN_ENV_VAR = "METORO_AUTH_TOKEN"

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

type Log struct {
	// The time that the log line was emitted in milliseconds since the epoch
	Time int64 `json:"time"`
	// The severity of the log line
	Severity string `json:"severity"`
	// The log message
	Message string `json:"message"`
	// The attributes of the log line
	LogAttributes map[string]string `json:"logAttributes"`
	// The attributes of the resource that emitted the log line
	ResourceAttributes map[string]string `json:"resourceAttributes"`
	// Service name
	ServiceName string `json:"serviceName"`
	// Environment
	Environment string `json:"environment"`
}

type GetLogsResponse struct {
	// The logs that match the filters
	Logs []Log `json:"logs"`
}

type GetTracesRequest struct {
	ServiceNames   []string            `json:"serviceNames"`
	StartTime     int64               `json:"startTime"`
	EndTime       int64               `json:"endTime"`
	Filters       map[string][]string `json:"filters"`
	ExcludeFilters map[string][]string `json:"excludeFilters"`
	PrevEndTime   *int64              `json:"prevEndTime"`
	Regexes       []string            `json:"regexes"`
	ExcludeRegexes []string           `json:"excludeRegexes"`
	Ascending     bool                `json:"ascending"`
	Environments  []string            `json:"environments"`
}