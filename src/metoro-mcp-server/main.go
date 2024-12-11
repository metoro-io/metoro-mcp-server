package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"os"
)

type MyTools struct {
	Name        string
	Description string
	Handler     any
}

var newTools = []MyTools{
	{
		Name:        "get_environments",
		Description: "Get Kubernetes environments/clusters, monitored by Metoro",
		Handler:     getEnvironmentsHandler,
	},
	{
		Name:        "get_services",
		Description: "Get services running in your Kubernetes cluster, monitored by Metoro",
		Handler:     getServicesHandler,
	},
	{
		Name:        "get_namespaces",
		Description: "Get namespaces in your Kubernetes cluster, monitored by Metoro",
		Handler:     getNamespacesHandler,
	},
	{
		Name:        "get_logs",
		Description: "Get logs from all/any services/hosts/pods running in your Kubernetes cluster in the last 5 minutes, monitored by Metoro",
		Handler:     getLogsHandler,
	},
	{
		Name:        "get_traces",
		Description: "Get traces from services running in your Kubernetes cluster in the last 5 minutes, monitored by Metoro",
		Handler:     getTracesHandler,
	},
	{
		Name:        "get_metric",
		Description: "Get metrics from your Kubernetes cluster, monitored by Metoro",
		Handler:     getMetricHandler,
	},
	//{
	//	Name:        "get_trace_metric",
	//	Description: "Get trace metrics from your Kubernetes cluster, monitored by Metoro",
	//	Handler:     getTraceMetricHandler,
	//},
	//{
	//	Name:        "get_trace_attributes",
	//	Description: "Get trace attributes from your Kubernetes cluster",
	//	Handler:     getTraceAttributesHandler,
	//},
	//{
	//	Name:        "get_trace_attribute_values_for_individual_attribute",
	//	Description: "Get trace attribute values for a specific attribute",
	//	Handler:     getTraceAttributeValuesForIndividualAttributeHandler,
	//},
	//{
	//	Name:        "get_profiles",
	//	Description: "Get profiles from your Kubernetes cluster",
	//	Handler:     getProfilesHandler,
	//},
	//{
	//	Name:        "get_k8s_events",
	//	Description: "Get Kubernetes events from your clusters with filtering options",
	//	Handler:     getK8sEventsHandler,
	//},
	//{
	//	Name:        "get_k8s_events_attributes",
	//	Description: "Get Kubernetes events attributes",
	//	Handler:     getK8sEventsAttributesHandler,
	//},
	//{
	//	Name:        "get_k8s_event_attribute_values_for_individual_attribute",
	//	Description: "Get Kubernetes event attribute values for a specific attribute",
	//	Handler:     getK8sEventAttributeValuesForIndividualAttributeHandler,
	//},
	//{
	//	Name:        "get_k8s_events_volume",
	//	Description: "Get Kubernetes events volume",
	//	Handler:     getK8sEventsVolumeHandler,
	//},
	//{
	//	Name:        "get_metricAttributes",
	//	Description: "Get metric attributes",
	//	Handler:     getMetricAttributesHandler,
	//},
	//{
	//	Name:        "get_metric_names",
	//	Description: "Get metric names",
	//	Handler:     getMetricNamesHandler,
	//},
	//{
	//	Name:        "get_metric_metadata",
	//	Description: "Get metric metadata",
	//	Handler:     getMetricMetadata,
	//},
	//{
	//	Name:        "get_pods",
	//	Description: "Get pods information from your Kubernetes cluster",
	//	Handler:     getPodsHandler,
	//},
	//{
	//	Name:        "get_k8s_service_information",
	//	Description: "Get detailed information about a Kubernetes service including its type (Deployment, DaemonSet, etc.), YAML configuration, and current running replicas (excluding HPA)",
	//	Handler:     getK8sServiceInformationHandler,
	//},
	//{
	//	Name:        "get_log_attributes",
	//	Description: "Get log attributes",
	//	Handler:     getLogAttributesHandler,
	//},
	//{
	//	Name:        "get_log_attribute_values_for_individual_attribute",
	//	Description: "Get log attribute values for a specific attribute",
	//	Handler:     getLogAttributeValuesForIndividualAttributeHandler,
	//},
	//{
	//	Name:        "get_nodes",
	//	Description: "Get nodes information from your Kubernetes cluster",
	//	Handler:     getNodesHandler,
	//},
	//{
	//	Name:        "get_node_info",
	//	Description: "Get detailed node information from your Kubernetes cluster",
	//	Handler:     getNodeInfoHandler,
	//},
	//{
	//	Name:        "get_service_summaries",
	//	Description: "Get service summaries from your Kubernetes cluster",
	//	Handler:     getServiceSummariesHandler,
	//},
	//{
	//	Name:        "get_alerts",
	//	Description: "Get alerts from your Kubernetes cluster",
	//	Handler:     getAlertsHandler,
	//},
	//{
	//	Name:        "get_alert_fires",
	//	Description: "Get alert fires from your Kubernetes cluster",
	//	Handler:     getAlertFiresHandler,
	//},
}

func main() {
	// Check if the appropriate environment variables are set
	if err := checkEnvVars(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

	done := make(chan struct{})

	mcpServer := mcpgolang.NewServer(stdio.NewStdioServerTransport())
	for _, tool := range newTools {
		err := mcpServer.RegisterTool(tool.Name, tool.Description, tool.Handler)
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
	if os.Getenv(METORO_API_URL_ENV_VAR) == "" {
		return fmt.Errorf("%s environment variable not set", METORO_API_URL_ENV_VAR)
	}
	if os.Getenv(METORO_AUTH_TOKEN_ENV_VAR) == "" {
		return fmt.Errorf("%s environment variable not set", METORO_AUTH_TOKEN_ENV_VAR)
	}
	return nil
}
