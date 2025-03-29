package tools

type MetoroTools struct {
	Name        string
	Description string
	Handler     any
}

var MetoroToolsList = []MetoroTools{
	{
		Name:        "get_environments",
		Description: "Get Kubernetes environments/clusters. This tool is useful for listing the kubernetes environments/clusters that are monitored by Metoro.",
		Handler:     GetEnvironmentsHandler,
	},
	{
		Name:        "get_services",
		Description: "Get services running in your Kubernetes cluster. Metoro treats the following Kubernetes resources as a 'service': Deployment, StatefulSet, DaemonSet. This tool is useful for listing the services/workloads running in your Kubernetes cluster.",
		Handler:     GetServicesHandler,
	},
	{
		Name:        "get_namespaces",
		Description: "Get namespaces in your Kubernetes cluster. This tool is useful for listing the namespaces in your Kubernetes cluster.",
		Handler:     GetNamespacesHandler,
	},
	{
		Name: "get_logs",
		Description: `Get logs from all or specific services/hosts/pods running in your Kubernetes cluster. Results are limited to 100 logs lines. How to use this tool:
                      First, use get_log_attributes tool to retrieve the available log attribute keys which can be used as Filter or ExcludeFilter keys for this tool.
                      Then use get_log_attribute_values_for_individual_attribute tool to get the possible values a log attribute key can be used for filtering logs.
                      Then you can use this tool (get_logs) to get the specific logs you are looking for.
                      e.g. Filter use case: get_logs with filters: {key: [value]} for including specific logs. Where the key was retrieved from get_log_attributes tool and the value was retrieved from get_log_attribute_values_for_individual_attribute tool. Multiple values for a key are ORed together.
                      If you want to filter log messages with a specific substring, you can use Regexes or ExcludeRegexes argument.`,
		Handler: GetLogsHandler,
	},
	{
		Name: "get_traces",
		Description: `Get traces from all or specific services/hosts/pods running in your Kubernetes cluster. Results are limited to 100 traces. How to use this tool:
					  First, use get_trace_attributes tool to retrieve the available trace attribute keys which can be used as Filter/ExcludeFilter keys.
                      Then use get_trace_attribute_values_for_individual_attribute tool to get the possible values a trace attribute key can be for filtering traces.
                      Then you can use this tool (get_traces) to get the specific traces you are looking for. 
                      e.g. Filter use case: get_traces with filters: {key: [value]} for including specific traces. Where key was retrieved from get_trace_attributes tool and value was retrieved from get_trace_attribute_values_for_individual_attribute tool. Multiple values for a key are ORed together. 
					  Regexes and ExcludeRegexes arguments can be used to filter traces endpoints that match the given regexes.`,
		Handler: GetTracesHandler,
	},
	{
		Name: "get_timeseries_data",
		Description: `Get one or more timeseries data for a metric or traces or logs or kubernetes resources. This tool is useful for understanding how the underlying type of data (specific/metric/trace/kubernetes resources/logs) change over time. You can also apply formulas to combine timeseries to calculate rates or ratios or differences etc. How to use this tool:
					  First you need the type of timeseries data you are requesting for. This can be one of metric or traces or logs or kubernetes resources. If it is metrics then you need to call the get_metric_names tool to get the available metric names which can be used as MetricName argument for this tool.
					  Then use get_attribute_keys tool to retrieve the available attribute keys and get_attribute_values to retrieve values for the attribute key that you are interested in to use in Filter/ExcludeFilter keys or Splits argument for this tool.
					  You can also use Splits argument to group the metric data by the given metric attribute keys. Only use the attribute keys and values that are available for the MetricName that are returned from get_attribute_keys and get_attribute_values tools.`,
		Handler: GetMultiMetricHandler,
	},
	{
		Name: "get_attribute_keys",
		Description: `Get the possible attribute keys for a specific type of data. This tool is useful for understanding the possible attribute keys that can be used for filtering the data. How to use this tool:
					  First you need the type of data you are requesting for. This can be one of metric or traces or logs or kubernetes resources.
					  Then you can call this tool to get the possible attribute keys for the given type of data.`,
		Handler: GetAttributeKeysHandler,
	},
	{
		Name: "get_attribute_values",
		Description: `"Get the possible values of an attribute key for a given type of data which can be one of metric trace logs or kubernetes_resource. This can be used as a value for a filtering key for filtering data. How to use this tool:
					  First you need the type of data you are requesting for. This can be one of metric or traces or logs or kubernetes resources. Then you need the attribute keys for the given type of data. You can use get_attribute_keys tool to get the available attribute keys for the given type of data.
					  Then you can call this tool to get the possible values for a given attribute key for the given type of data. If you want to get the possible values for a metric attribute key you can use the get_metric_names tool to get the available metric names which can be used as MetricName argument for this tool and then use get_attribute_keys tool to get the available attribute keys and get_attribute_values to get values for the key which can be used as Filter/ExcludeFilter keys for`,
		Handler: GetAttributeValuesHandler,
	},
	//{
	//	Name:        "get_trace_attribute_values_for_individual_attribute",
	//	Description: "Get trace the possible values a trace attribute key can be used as a value for filtering traces",
	//	Handler:     GetTraceAttributeValuesForIndividualAttributeHandler,
	//},
	{
		Name:        "get_profiles",
		Description: "Get profiles of your services running in your Kubernetes cluster. This tool is useful for answering performance related questions for a specific service. It provides information about which functions taking time in the service.",
		Handler:     GetProfilesHandler,
	},
	{
		Name: "get_k8s_events",
		Description: `Get the Kubernetes events from your clusters. Kubernetes events are useful for understanding what is happening in your cluster. 
They are emitted by the Kubernetes API server when there is a change in the state of the cluster. How to use this tool:
First, use get_k8s_events_attributes tool to retrieve the available Kubernetes event attribute keys which can be used as Filter/ExcludeFilter keys for this tool.
Then use get_k8s_event_attribute_values_for_individual_attribute tool to get the possible values a Kubernetes event attribute key can be for filtering Kubernetes events.
And then you can call this tool (get_k8s_events) to get the specific events you are looking for. e.g. Filter use case: get_k8s_events with filters: {key: [value]} for including specific Kubernetes events.`,
		Handler: GetK8sEventsHandler,
	},
	{
		Name:        "get_k8s_events_attributes",
		Description: "Get possible attribute keys for Kubernetes events which can be used for filtering them.",
		Handler:     GetK8sEventsAttributesHandler,
	},
	{
		Name:        "get_k8s_event_attribute_values_for_individual_attribute",
		Description: "Get possible attribute values for a specific Kubernetes event attribute key. E.g. EventType attribute key might have values like Normal, Warning, etc.",
		Handler:     GetK8sEventAttributeValuesForIndividualAttributeHandler,
	},
	{
		Name:        "get_k8s_events_volume",
		Description: "Get the timeseries data for the number of Kubernetes events in your cluster, whether its filtered by a specific attribute or not. The volume of events are split by EventType so you can see the breakdown of Warning/Normal events.",
		Handler:     GetK8sEventsVolumeHandler,
	},
	{
		Name:        "get_metric_names",
		Description: "Get available metric names to query. These metric names can be used as MetricName argument for get_metric, get_metric_metadata and get_timeseries_data and get_attribute_keys tools.",
		Handler:     GetMetricNamesHandler,
	},
	{
		Name:        "get_metric_metadata",
		Description: "Get metric description and type and unit for a metric. This tool can be used to get detailed information about a metric including its type, unit and description. Use this tool after getting the metric name that you are interested in from the get_metric_names tool and before calling the get_timeseries_data tool to understand the metric better.",
		Handler:     GetMetricMetadata,
	},
	{
		Name:        "get_pods",
		Description: "Get the pods that are running in your cluster. You must provide either a ServiceName to get pods for a specific service or a NodeName to get pods running on a specific node.",
		Handler:     GetPodsHandler,
	},
	{
		Name:        "get_k8s_service_information",
		Description: "Get detailed information including the YAML of a Kubernetes service. This tool is useful for understanding the configuration of a service.",
		Handler:     GetK8sServiceInformationHandler,
	},
	//{
	//	Name:        "get_log_attribute_values_for_individual_attribute",
	//	Description: "Get possible values for a specific log attribute key which can be used for filtering logs.",
	//	Handler:     GetLogAttributeValuesForIndividualAttributeHandler,
	//},
	{
		Name:        "get_nodes",
		Description: "Get the nodes that are running in your cluster. To use this tool, first call get_node_attributes to get the possible node attribute keys and values which can be used for filtering nodes.",
		Handler:     GetNodesHandler,
	},
	{
		Name:        "get_node_attributes",
		Description: "Get possible node attribute keys and values which can be used for filtering nodes.",
		Handler:     GetNodeAttributesHandler,
	},
	{
		Name:        "get_node_info",
		Description: "Get detailed node information about a specific node. This tool provides information about the node's capacity, allocatable resources, and usage, yaml, node type, OS and Kernel information.",
		Handler:     GetNodeInfoHandler,
	},
	{
		Name:        "get_service_summaries",
		Description: "Get summaries of services/workloads running in your Kubernetes cluster. The summary includes the number of requests, errors (5xx and 4xx), P50, p95, p99 latencies. This tool is useful for understanding the performance of your services at a high level for a given relative or abosulute time range.",
		Handler:     GetServiceSummariesHandler,
	},
	{
		Name:        "get_alerts",
		Description: "Get list of alerts from your Kubernetes cluster. These alerts are configured by the user in Metoro therefore it may not have full coverage for all the issues that might occur in the cluster.",
		Handler:     GetAlertsHandler,
	},
	{
		Name:        "get_alert_fires",
		Description: "Get list of alert fires from your Kubernetes cluster. Alert fires are the instances when an alert is triggered. This tool provides information about the alert name, the time it was triggered, the time it recovered, the environment, and the service name (if available) and the alert trigger message.",
		Handler:     GetAlertFiresHandler,
	},
	{
		Name: "create_dashboard",
		Description: `Create a dashboard with the described metrics. This tool is useful for creating a dashboard with the metrics you are interested in.
											  How to use this tool:
					  First use get_metric_names tool to retrieve the available metric names which can be used as MetricName argument for this tool and then use get_attribute_keys tool to retrieve the available attribute keys and get_attribute_values for getting the values for the attribute key that you are interested in to use in Filter/ExcludeFilter keys or Splits argument for MetricChartWidget argument for this tool.
					  You can also use Splits argument to group the metric data by the given metric attribute keys. Only use the attribute keys and values that are available for the MetricName that are returned from get_attribute_keys and get_attribute_values tools.`,
		Handler: CreateDashboardHandler,
	},
	{
		Name:        "get_source_repository",
		Description: "Get the source repository URL/path for a specific service. This tool is useful for finding where the code for a service is stored. You need to provide the service name time range and optionally specific environments to search in.",
		Handler:     GetSourceRepositoryHandler,
	},
}
