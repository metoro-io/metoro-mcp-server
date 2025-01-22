package resources

type MetoroResource struct {
	Path        string
	Name        string
	Description string
	ContentType string
	Handler     any
}

var MetoroResourcesList = []MetoroResource{
	{
		Path:        "api://environments",
		Name:        "environments",
		Description: "This resource provides a list of names of the kubernetes clusters/environments monitored by Metoro",
		ContentType: "text/plain",
		Handler:     EnvironmentResourceHandler,
	},
	{
		Path:        "api://namespaces",
		Name:        "namespaces",
		Description: "This resource provides a list of namespaces in the kubernetes clusters/environments monitored by Metoro",
		ContentType: "text/plain",
		Handler:     NamespacesResourceHandler,
	},
	{
		Path:        "api://services",
		Name:        "services",
		Description: "This resource provides a list of services running in the kubernetes clusters/environments monitored by Metoro",
		ContentType: "text/plain",
		Handler:     ServicesResourceHandler,
	},
	{
		Path:        "api://traceAttributes",
		Name:        "traceAttributes",
		Description: "Provides a list of trace attribute keys that are available to be used for filtering or grouping traces. These trace attribute keys should be used as Filter/ExcludeFilter keys or Splits for get_traces, get_trace_metric and get_trace_attribute_values_for_individual_attribute tools arguments.",
		ContentType: "text/plain",
		Handler:     TraceAttributesResourceHandler,
	},
	{
		Path:        "api://k8sEventAttributes",
		Name:        "k8sEventAttributes",
		Description: "Provides a list of Kubernetes Event's attribute keys that are available to be used for filtering or grouping K8s Events. These K8s Event attribute keys should be used as Filter/ExcludeFilter keys or Splits for get_k8s_events, get_k8s_events_volume and get_k8s_events_volume tools arguments.",
		ContentType: "text/plain",
		Handler:     K8sEventsAttributesResourceHandler,
	},
	{
		Path:        "api://metrics",
		Name:        "metricNames",
		Description: "Provides a list of available metric names that can be used for as MetricName arguments to get_metric, get_metric_metadata and get_metric_attributes tools to get metrics data.",
		ContentType: "text/plain",
		Handler:     MetricsResourceHandler,
	},
	{
		Path:        "api://logAttributes",
		Name:        "logAttributes",
		Description: "Provides a list of log attribute keys that are available to be used for filtering or grouping logs. These log attribute keys should be used as Filter/ExcludeFilter keys or Splits for get_logs, get_log_attribute_values_for_individual_attribute tools arguments.",
		ContentType: "text/plain",
		Handler:     LogAttributesResourceHandler,
	},
	{
		Path:        "api://nodes",
		Name:        "nodes",
		Description: "Provides a list of nodes in the kubernetes clusters/environments monitored by Metoro. Any of these nodes/instances can be used as a filter/exclude for get_metric tool with the key 'kubernetes.io/hostname' and value as the node names in this resource.",
		ContentType: "text/plain",
		Handler:     NodesResourceHandler,
	},
}
