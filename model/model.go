package model

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
	// The type of the function
	FunctionType FunctionType `json:"functionType" jsonschema:"required,enum=monotonicDifference,enum=valueDifference,enum=perSecond,description=The type of the function to apply to the metric. Do not guess the function type. Use the available ones: perSecond or valueDifference or monotonicDifference."`
	//// The payload of the function
	//// TODO: If we have more payloads this can be an interface but for now its a math expression since its the only payload.
	//FunctionPayload MathExpression `json:"functionPayload" jsonschema:"description=The payload of the customMathExpression. this is only set for customMathExpression. "`
}

type MathExpression struct {
	Variables  []string `json:"variables" jsonschema:"description=The variables to use in the math expression. For now this should always be ['a'] if set"`
	Expression string   `json:"expression" jsonschema:"description=The math expression to apply to the metric. For example if you want to divide the metric by 60 you would set the expression as a / 60"`
}

type FunctionType string

const (
	MonotonicDifference FunctionType = "monotonicDifference"
	ValueDifference     FunctionType = "valueDifference"
)

type GetMetricRequest struct {
	// MetricName is the name of the metric to get
	MetricName string `json:"metricName" jsonschema:"required,description=Name of the metric to get the timeseries data for. Do not guess the metricName, get the possible values from get_metric_names tool"`
	// Required: Start time of when to get the logs in seconds since epoch
	StartTime int64 `json:"startTime" jsonschema:"required,description=Start time of when to get the metrics in seconds since epoch"`
	// Required: End time of when to get the logs in seconds since epoch
	EndTime int64 `json:"endTime" jsonschema:"required,description=Start time of when to get the metrics in seconds since epoch"`
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

type MetricAttributesRequest struct {
	StartTime        int64               `json:"startTime"`
	EndTime          int64               `json:"endTime"`
	MetricName       string              `json:"metricName"`
	FilterAttributes map[string][]string `json:"filterAttributes"`
}

type FuzzyMetricsRequest struct {
	MetricFuzzyMatch string   `json:"metricFuzzyMatch"`
	Environments     []string `json:"environments"`
	StartTime        int64    `json:"startTime"`
	EndTime          int64    `json:"endTime"`
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

type GetSingleTraceSummaryRequest struct {
	TracesSummaryRequest
	// The attribute to get the summary for
	Attribute string `json:"attribute"`
}

type TracesSummaryRequest struct {
	// Required: Start time of when to get the service summaries in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds since epoch
	EndTime int64 `json:"endTime"`

	// The filters to apply to the trace summary, so for example, if you want to get traces for a specific service
	// you can pass in a filter like {"service_name": ["microservice_a"]}
	Filters map[string][]string `json:"filters"`
	// ExcludeFilters are used to exclude traces based on a filter
	ExcludeFilters map[string][]string `json:"excludeFilters"`

	// Regexes are used to filter traces based on a regex inclusively
	Regexes []string `json:"regexes"`
	// ExcludeRegexes are used to filter traces based on a regex exclusively
	ExcludeRegexes []string `json:"excludeRegexes"`

	// Optional: The name of the service to get the trace metrics for
	// Acts as an additional filter
	ServiceNames []string `json:"serviceNames"`

	// Environments is the environments to get the traces for. If empty, all environments will be included
	Environments []string `json:"environments"`
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

type GetServiceSummariesRequest struct {
	// Required: Start time of when to get the service summaries in seconds
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds
	EndTime int64 `json:"endTime"`
	// If empty, all services across all environments will be returned
	Environments []string `json:"environments"`
	// Required: The namespace of the services to get summaries for. If empty, return services from all namespaces
	Namespace string `json:"namespace"`
}

// Dasboarding structs
type SetDashboardRequest struct {
	Name             string `json:"name"`
	Id               string `json:"id"`
	DashboardJson    string `json:"dashboardJson"`
	DefaultTimeRange string `json:"defaultTimeRange"`
}

// WidgetType is an enum representing different types of widgets
type WidgetType string

const (
	MetricChartWidgetType WidgetType = "MetricChart"
	GroupWidgetType       WidgetType = "Group"
	MarkdownWidgetType    WidgetType = "Markdown"
)

// WidgetPosition represents the position of a widget relative to its parent
type WidgetPosition struct {
	X *int `json:"x,omitempty"`
	Y *int `json:"y,omitempty"`
	W *int `json:"w,omitempty" jsonschema:"required,description=The width of the widget. The dashboard is divided into 12 columns.For example a sensible value for a graph would be 6"`
	H *int `json:"h,omitempty" jsonschema:"required,description=The height of the widget. Each row is 128px. A sensible value for a graph would be 3."`
}

// Widget is the base interface for all widget types
type Widget struct {
	WidgetType WidgetType      `json:"widgetType" jsonschema:"required,description=The type of the widget. This can be MetricChart / Group / Markdown"`
	Position   *WidgetPosition `json:"position,omitempty" jsonschema:"description=The position of the widget in the dashboard"`
}

//// Variable represents a configurable variable in a widget
//type Variable struct {
//	Name              string     `json:"name"`
//	Key               string     `json:"key"`
//	DefaultValue      string     `json:"defaultValue"`
//	DefaultType       MetricType `json:"defaultType"`
//	DefaultMetricName *string    `json:"defaultMetricName,omitempty"`
//	IsOverrideable    bool       `json:"isOverrideable"`
//}

// GroupWidget represents a group of widgets
type GroupWidget struct {
	Widget   `json:",inline"`
	Title    *string             `json:"title,omitempty" jsonschema:"description=The title of the group widget if present"`
	Children []MetricChartWidget `json:"children" jsonschema:"description=The children widgets of the group widget. The children are MetricChartWidgets."`
	//Variables []Variable `json:"variables,omitempty"`
}

// MetricChartWidget represents a metric chart widget
type MetricChartWidget struct {
	Widget         `json:",inline"`
	MetricName     string              `json:"metricName" jsonschema:"description=The name of the metric to use in the chart if MetricType is metric. If MetricType is trace, this is not used and can be empty. This value is same as the metricName in the get_metric tool and the possible metricNames can be found in the get_metric_names tool"`
	Filters        map[string][]string `json:"filters,omitempty" jsonschema:"description=The filters to apply to the metric. This is the same as the filters in the get_metric or get_trace_metric tool depending on the MetricType"`
	ExcludeFilters map[string][]string `json:"excludeFilters,omitempty" jsonschema:"description=The exclude filters to apply to the metric. This is the same as the exclude filters in the get_metric or get_trace_metric tool depending on the MetricType"`
	Splits         []string            `json:"splits,omitempty" jsonshcema:"description=Splits will allow you to group/split metrics by an attribute. This is useful if you would like to see the breakdown of a particular metric by an attribute. For example if you want to see the breakdown of the metric by service_name you would set the splits as ['service_name']"`
	Aggregation    string              `json:"aggregation" jsonschema:"description=The aggregation to apply to the metrics. This is the same as the aggregation in the get_metric or get_trace_metric tool depending on the MetricType"`
	Title          *string             `json:"title,omitempty" jsonschema:"description=The title of the metric chart widget if present"`
	Type           ChartType           `json:"type" jsonschema:"description=The type of the chart to display. Possible values are line / bar."`
	MetricType     MetricType          `json:"metricType" jsonschema:"description=The type of the metric to use in the chart. Possible values are metric / trace. If metric, the metricName should be used."`
	Functions      []MetricFunction    `json:"functions" jsonschema:"description=The functions to apply to the metric. This is the same as the functions in the get_metric or get_trace_metric tool depending on the MetricType"`
}

// MarkdownWidget represents a markdown content widget
type MarkdownWidget struct {
	Widget  `json:",inline"`
	Content string `json:"content"`
}
type ChartType string

const (
	ChartTypeLine ChartType = "line"
	ChartTypeBar  ChartType = "bar"
)

type MetricType string

const (
	Metric MetricType = "metric" // please excuse the bad naming... this is a metric timeseries type.
	Trace  MetricType = "trace"  // trace timeseries type.

	Logs MetricType = "logs" // log timeseries type.

	KubernetesResource MetricType = "kubernetes_resource" // kubernetes resource timeseries type.
)

type GetLogMetricRequest struct {
	GetLogsRequest
	Splits     []string         `json:"splits" jsonschema:"description=Splits will allow you to group/split metrics by an attribute. This is useful if you would like to see the breakdown of a particular metric by an attribute. For example if you want to see the breakdown of the metric by service.name you would set the splits as ['service.name']"`
	Functions  []MetricFunction `json:"functions" jsonschema:"description=The functions to apply to the log metric. Available functions are monotonicDifference which will calculate the difference between the current and previous value of the metric (negative values will be set to 0) and valueDifference which will calculate the difference between the current and previous value of the metric or MathExpression e.g. a / 60"`
	BucketSize int64            `json:"bucketSize" jsonschema:"description=The size of each datapoint bucket in seconds if not provided metoro will select the best bucket size for the given duration for performance and clarity"`
}

type GetKubernetesResourceRequest struct {
	// Required: Start time of when to get the logs in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the logs in seconds since epoch
	EndTime int64 `json:"endTime"`
	//
	//// Optional: The name of the service to get the k8s events metrics for
	//// Acts as an additional filter
	//ServiceNames []string `json:"serviceNames"`

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

	// Environments is a list of environments to filter the k8s events by. If empty, all environments will be included
	Environments []string `json:"environments"`
}

type GetMultiMetricRequest struct {
	// Required: Start time of when to get the service summaries in seconds
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds
	EndTime  int64                 `json:"endTime"`
	Metrics  []SingleMetricRequest `json:"metrics" jsonschema:"required,description=Array of metrics to get the timeseries data for"`
	Formulas []Formula             `json:"formulas" jsonschema:"description=Optional formulas to combine metrics/log metrics/trace metrics. Formula should only consist of formulaIdentifier of the metrics/logs/traces in the metrics array"`
}

type SingleMetricRequest struct {
	Type   string                 `json:"type" jsonschema:"required,enum=metric,enum=trace,enum=logs,enum=kubernetes_resource,description=Type of metric to retrieve"`
	Metric *GetMetricRequest      `json:"metric,omitempty" jsonschema:"description=Metric request details when type is 'metric'"`
	Trace  *GetTraceMetricRequest `json:"trace,omitempty" jsonschema:"description=Trace metric request details when type is 'trace'"`
	Logs   *GetLogMetricRequest   `json:"logs,omitempty" jsonschema:"description=Log metric request details when type is 'logs'"`
	// TODO: Add kubernetes resource request
	ShouldNotReturn   bool   `json:"shouldNotReturn" jsonschema:"description=If true result won't be returned (useful for formulas)"`
	FormulaIdentifier string `json:"formulaIdentifier" jsonschema:"description=Identifier to reference this metric in formulas"`
}

type Formula struct {
	Formula string `json:"formula" jsonschema:"description=Math expression combining metric results using their formula identifiers"`
}

type GetMetricAttributesRequest struct {
	// Required: The metric name to get the summary for
	MetricName string `json:"metricName"`
	// Required: Start time of when to get the service summaries in seconds since epoch
	StartTime int64 `json:"startTime"`
	// Required: End time of when to get the service summaries in seconds since epoch
	EndTime int64 `json:"endTime"`
	// Environments is the environments to get the traces for. If empty, all environments will be included
	Environments []string `json:"environments"`
}

type MultiMetricAttributeKeysRequest struct {
	Type   string                      `json:"type"`
	Metric *GetMetricAttributesRequest `json:"metric,omitempty"`
	// Currently trace and logs and kubernetes resource do not have any request parameters
	// Only metric has request parameters
}

type GetAttributeValuesRequest struct {
	Type       MetricType                    `json:"type"`
	Attribute  string                        `json:"attribute"`
	Limit      *int                          `json:"limit"`
	Metric     *GetMetricAttributesRequest   `json:"metric,omitempty"`
	Trace      *TracesSummaryRequest         `json:"trace,omitempty"`
	Logs       *LogSummaryRequest            `json:"logs,omitempty"`
	Kubernetes *GetKubernetesResourceRequest `json:"kubernetes,omitempty"`
}

type GetAttributeKeysResponse struct {
	// The attribute values
	Attributes []string `json:"attributes"`
}
