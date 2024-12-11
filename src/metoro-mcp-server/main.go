package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/resources"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/tools"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
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
		Description: "Get Kubernetes environments/clusters, monitored by Metoro",
		Handler:     tools.GetEnvironmentsHandler,
	},
	{
		Name:        "get_services",
		Description: "Get services running in your Kubernetes cluster, monitored by Metoro",
		Handler:     tools.GetServicesHandler,
	},
	{
		Name:        "get_namespaces",
		Description: "Get namespaces in your Kubernetes cluster, monitored by Metoro",
		Handler:     tools.GetNamespacesHandler,
	},
	{
		Name:        "get_logs",
		Description: "Get logs from all/any services/hosts/pods running in your Kubernetes cluster in the last 5 minutes, monitored by Metoro",
		Handler:     tools.GetLogsHandler,
	},
	{
		Name:        "get_traces",
		Description: "Get traces from services running in your Kubernetes cluster in the last 5 minutes, monitored by Metoro",
		Handler:     tools.GetTracesHandler,
	},
	{
		Name:        "get_metric",
		Description: "Get metrics from your Kubernetes cluster, monitored by Metoro",
		Handler:     tools.GetMetricHandler,
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

func main() {
	// Check if the appropriate environment variables are set
	if err := checkEnvVars(); err != nil {
		panic(err)
	}

	done := make(chan struct{})

	mcpServer := mcpgolang.NewServer(stdio.NewStdioServerTransport())
	for _, tool := range metoroTools {
		err := mcpServer.RegisterTool(tool.Name, tool.Description, tool.Handler)
		if err != nil {
			panic(err)
		}
	}

	err := mcpServer.RegisterResource(
		"api://environments",
		"environments",
		"This resource provides a list of names of the kubernetes clusters/environments monitored by Metoro",
		"text/plain",
		resources.EnvironmentResourceHandler)

	if err != nil {
		panic(err)
	}

	err = mcpServer.Serve()
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
