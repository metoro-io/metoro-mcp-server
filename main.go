package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/metoro-io/metoro-mcp-server/resources"
	"github.com/metoro-io/metoro-mcp-server/tools"
	"github.com/metoro-io/metoro-mcp-server/utils"
	"os"
)

type MetoroTools struct {
	Name        string
	Description string
	Handler     any
}

var metoroTools = []MetoroTools{
	{
		Name:        "get_environments",
		Description: "Get Kubernetes environments/clusters",
		Handler:     tools.GetEnvironmentsHandler,
	},
	{
		Name:        "get_services",
		Description: "Get services running in your Kubernetes cluster. Metoro treats the following Kubernetes resources as a 'service': Deployment, StatefulSet, DaemonSet",
		Handler:     tools.GetServicesHandler,
	},
	{
		Name:        "get_namespaces",
		Description: "Get namespaces in your Kubernetes cluster",
		Handler:     tools.GetNamespacesHandler,
	},
	{
		Name: "get_logs",
		Description: `Get logs from all or specific services/hosts/pods running in your Kubernetes cluster. Results are limited to 100 logs lines. How to use this tool:
                      First, use get_log_attributes tool to retrieve the available log attribute keys which can be used as Filter or ExcludeFilter keys for this tool.
                      e.g. Filter use case: get_logs with filters: {key: [value]} for including specific logs. Where key was retrieved from get_log_attributes tool.
                      If you want to filter log messages with a specific substring, you can use Regexes or ExcludeRegexes argument.`,
		Handler: tools.GetLogsHandler,
	},
	{
		Name: "get_traces",
		Description: `Get traces from all or specific services/hosts/pods running in your Kubernetes cluster. Results are limited to 100 traces. How to use this tool:
					  First, use get_trace_attributes tool to retrieve the available trace attribute keys which can be used as Filter/ExcludeFilter keys.
                      e.g. Filter use case: get_traces with filters: {key: [value]} for including specific traces. Where key was retrieved from get_trace_attributes tool.
					  Regexes and ExcludeRegexes arguments can be used to filter traces endpoints that match the given regexes.`,
		Handler: tools.GetTracesHandler,
	},
	{
		Name: "get_metric",
		Description: `Get a specific metric's timeseries data. How to use this tool: 
                      First, use get_metric_names tool to retrieve the available metric names which can be used as MetricName argument for this tool.
                      Then use get_metric_attributes tool to retrieve the available attribute keys for a specific MetricName which can be used as Filter/ExcludeFilter keys for this tool.
                      You can also use Splits argument to group the metric data by the given metric attribute keys.`,
		Handler: tools.GetMetricHandler,
	},
	{
		Name:        "get_trace_metric",
		Description: "Get trace metrics from your Kubernetes cluster, monitored by Metoro",
		Handler:     tools.GetTraceMetricHandler,
	},
	{
		Name:        "get_trace_attributes",
		Description: "Get trace attributes from your Kubernetes cluster",
		Handler:     tools.GetTraceAttributesHandler,
	},
	{
		Name:        "get_trace_attribute_values_for_individual_attribute",
		Description: "Get trace attribute values for a specific attribute",
		Handler:     tools.GetTraceAttributeValuesForIndividualAttributeHandler,
	},
	{
		Name:        "get_profiles",
		Description: "Get profiles from your Kubernetes cluster",
		Handler:     tools.GetProfilesHandler,
	},
	{
		Name:        "get_k8s_events",
		Description: "Get Kubernetes events from your clusters with filtering options",
		Handler:     tools.GetK8sEventsHandler,
	},
	{
		Name:        "get_k8s_events_attributes",
		Description: "Get Kubernetes events attributes",
		Handler:     tools.GetK8sEventsAttributesHandler,
	},
	{
		Name:        "get_k8s_event_attribute_values_for_individual_attribute",
		Description: "Get Kubernetes event attribute values for a specific attribute",
		Handler:     tools.GetK8sEventAttributeValuesForIndividualAttributeHandler,
	},
	{
		Name:        "get_k8s_events_volume",
		Description: "Get Kubernetes events volume",
		Handler:     tools.GetK8sEventsVolumeHandler,
	},
	{
		Name:        "get_metricAttributes",
		Description: "Get metric attributes",
		Handler:     tools.GetMetricAttributesHandler,
	},
	{
		Name:        "get_metric_names",
		Description: "Get metric names",
		Handler:     tools.GetMetricNamesHandler,
	},
	{
		Name:        "get_metric_metadata",
		Description: "Get metric metadata",
		Handler:     tools.GetMetricMetadata,
	},
	{
		Name:        "get_pods",
		Description: "Get pods information from your Kubernetes cluster",
		Handler:     tools.GetPodsHandler,
	},
	{
		Name:        "get_k8s_service_information",
		Description: "Get detailed information about a Kubernetes service including its type (Deployment, DaemonSet, etc.), YAML configuration, and current running replicas (excluding HPA)",
		Handler:     tools.GetK8sServiceInformationHandler,
	},
	{
		Name:        "get_log_attributes",
		Description: "Get log attributes",
		Handler:     tools.GetLogAttributesHandler,
	},
	{
		Name:        "get_log_attribute_values_for_individual_attribute",
		Description: "Get log attribute values for a specific attribute",
		Handler:     tools.GetLogAttributeValuesForIndividualAttributeHandler,
	},
	{
		Name:        "get_nodes",
		Description: "Get nodes information from your Kubernetes cluster",
		Handler:     tools.GetNodesHandler,
	},
	{
		Name:        "get_node_info",
		Description: "Get detailed node information from your Kubernetes cluster",
		Handler:     tools.GetNodeInfoHandler,
	},
	{
		Name:        "get_service_summaries",
		Description: "Get service summaries from your Kubernetes cluster",
		Handler:     tools.GetServiceSummariesHandler,
	},
	{
		Name:        "get_alerts",
		Description: "Get alerts from your Kubernetes cluster",
		Handler:     tools.GetAlertsHandler,
	},
	{
		Name:        "get_alert_fires",
		Description: "Get alert fires from your Kubernetes cluster",
		Handler:     tools.GetAlertFiresHandler,
	},
}

type MetoroResource struct {
	Path        string
	Name        string
	Description string
	ContentType string
	Handler     any
}

var metoroResources = []MetoroResource{
	{
		Path:        "api://environments",
		Name:        "environments",
		Description: "This resource provides a list of names of the kubernetes clusters/environments monitored by Metoro",
		ContentType: "text/plain",
		Handler:     resources.EnvironmentResourceHandler,
	},
	{
		Path:        "api://namespaces",
		Name:        "namespaces",
		Description: "This resource provides a list of namespaces in the kubernetes clusters/environments monitored by Metoro",
		ContentType: "text/plain",
		Handler:     resources.NamespacesResourceHandler,
	},
	{
		Path:        "api://services",
		Name:        "services",
		Description: "This resource provides a list of services running in the kubernetes clusters/environments monitored by Metoro",
		ContentType: "text/plain",
		Handler:     resources.ServicesResourceHandler,
	},
	{
		Path:        "api://traceAttributes",
		Name:        "traceAttributes",
		Description: "Provides a list of trace attribute keys that are available to be used for filtering or grouping traces. These trace attribute keys should be used as Filter/ExcludeFilter keys or Splits for get_traces, get_trace_metric and get_trace_attribute_values_for_individual_attribute tools arguments.",
		ContentType: "text/plain",
		Handler:     resources.TraceAttributesResourceHandler,
	},
	{
		Path:        "api://k8sEventAttributes",
		Name:        "k8sEventAttributes",
		Description: "Provides a list of Kubernetes Event's attribute keys that are available to be used for filtering or grouping K8s Events. These K8s Event attribute keys should be used as Filter/ExcludeFilter keys or Splits for get_k8s_events, get_k8s_events_volume and get_k8s_events_volume tools arguments.",
		ContentType: "text/plain",
		Handler:     resources.K8sEventsAttributesResourceHandler,
	},
	{
		Path:        "api://metrics",
		Name:        "metricNames",
		Description: "Provides a list of available metric names that can be used for as MetricName arguments to get_metric, get_metric_metadata and get_metricAttributes tools to get metrics data.",
		ContentType: "text/plain",
		Handler:     resources.MetricsResourceHandler,
	},
	{
		Path:        "api://logAttributes",
		Name:        "logAttributes",
		Description: "Provides a list of log attribute keys that are available to be used for filtering or grouping logs. These log attribute keys should be used as Filter/ExcludeFilter keys or Splits for get_logs, get_log_attribute_values_for_individual_attribute tools arguments.",
		ContentType: "text/plain",
		Handler:     resources.LogAttributesResourceHandler,
	},
	{
		Path:        "api://nodes",
		Name:        "nodes",
		Description: "Provides a list of nodes in the kubernetes clusters/environments monitored by Metoro. Any of these nodes/instances can be used as a filter/exclude for get_metric tool with the key 'kubernetes.io/hostname' and value as the node names in this resource.",
		ContentType: "text/plain",
		Handler:     resources.NodesResourceHandler,
	},
}

func main() {
	// Check if the appropriate environment variables are set
	if err := checkEnvVars(); err != nil {
		panic(err)
	}

	done := make(chan struct{})

	mcpServer := mcpgolang.NewServer(stdio.NewStdioServerTransport())

	// Add tools
	for _, tool := range metoroTools {
		err := mcpServer.RegisterTool(tool.Name, tool.Description, tool.Handler)
		if err != nil {
			panic(err)
		}
	}

	// Add resources
	for _, resource := range metoroResources {
		err := mcpServer.RegisterResource(
			resource.Path,
			resource.Name,
			resource.Description,
			resource.ContentType,
			resource.Handler)
		if err != nil {
			panic(err)
		}
	}

	err := mcpServer.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}

func checkEnvVars() error {
	if os.Getenv(utils.METORO_API_URL_ENV_VAR) == "" {
		return fmt.Errorf("%s environment variable not set", utils.METORO_API_URL_ENV_VAR)
	}
	if os.Getenv(utils.METORO_AUTH_TOKEN_ENV_VAR) == "" {
		return fmt.Errorf("%s environment variable not set", utils.METORO_AUTH_TOKEN_ENV_VAR)
	}
	return nil
}
