package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"os"
)

var handlers = map[string]server.ToolHandlerFunc{
	"get_environments": getEnvironmentsHandler,
	"get_services":     getServicesHandler,
	"get_namespaces":   getNamespacesHandler,
	"get_logs":         getLogsHandler,
	"get_traces":       getTracesHandler,
	"get_metric":       getMetricHandler,
}

var tools = []mcp.Tool{
	mcp.NewTool("get_environments",
		mcp.WithDescription("Get Kubernetes environments/clusters, monitored by Metoro"),
	),
	mcp.NewTool("get_services",
		mcp.WithDescription("Get services running in your Kubernetes cluster, monitored by Metoro"),
	),
	mcp.NewTool("get_namespaces",
		mcp.WithDescription("Get namespaces in your Kubernetes cluster, monitored by Metoro"),
	),
	mcp.NewTool("get_logs",
		mcp.WithDescription("Get logs from all/any services/hosts/pods running in your Kubernetes cluster in the last 5 minutes, monitored by Metoro"),
		// TODO: Fix the issue with 5 minutes hardcoded
		mcp.WithString("filters",
			mcp.Description("The filters to apply to the logs. It is a stringified map[string]string[], e.g., '{\"service.name\": [\"/k8s/namespaceX/serviceX\"]}' should return logs for serviceX in namespaceX"),
		),
		mcp.WithString("excludeFilters",
			mcp.Description("The filters that should be excluded from the logs. It is a stringified map[string]string[] e.g., '{\"service.name\": [\"/k8s/namespaceX/serviceX\"]}' should return all logs except for serviceX in namespaceX"),
		),
		mcp.WithString("regexes",
			mcp.Description("JSON array of regexes as a string to filter logs based on a regex inclusively"),
		),
		mcp.WithString("excludeRegexes",
			mcp.Description("JSON array of regexes as a string to filter logs based on a regex exclusively"),
		),
		mcp.WithBoolean("ascending",
			mcp.Description("Whether to return logs in ascending order or not"),
		),
		mcp.WithString("environments",
			mcp.Description("JSON array of cluster/environments as a string. If empty, all clusters will be included"),
		),
	),
	mcp.NewTool("get_traces",
		mcp.WithDescription("Get traces from services running in your Kubernetes cluster in the last 5 minutes, monitored by Metoro"),
		mcp.WithString("serviceNames",
			mcp.Description("JSON array of service names as a string to filter traces for specific services"),
		),
		mcp.WithString("filters",
			mcp.Description("The filters to apply to the traces. It is a stringified map[string]string[], e.g., '{\"service.name\": [\"serviceX\"]}' should return traces for serviceX"),
		),
		mcp.WithString("excludeFilters",
			mcp.Description("The filters that should be excluded from the traces. It is a stringified map[string]string[] e.g., '{\"service.name\": [\"serviceX\"]}' should return all traces except for serviceX"),
		),
		mcp.WithString("regexes",
			mcp.Description("JSON array of regexes as a string to filter traces based on a regex inclusively"),
		),
		mcp.WithString("excludeRegexes",
			mcp.Description("JSON array of regexes as a string to filter traces based on a regex exclusively"),
		),
		mcp.WithBoolean("ascending",
			mcp.Description("Whether to return traces in ascending order or not"),
		),
		mcp.WithString("environments",
			mcp.Description("JSON array of cluster/environments as a string. If empty, all clusters will be included"),
		),
	),
	mcp.NewTool("get_metric",
		mcp.WithDescription("Get timeseries data of metrics from your Kubernetes cluster"),
		mcp.WithString("metricName",
			mcp.Description("The name of the metric to get"),
		),
		// TODO: Fix the issue with startTime and endTime
		//mcp.WithNumber("startTime",
		//	mcp.Description("Start time of when to get the metrics in seconds since epoch"),
		//),
		//mcp.WithNumber("endTime",
		//	mcp.Description("End time of when to get the metrics in seconds since epoch"),
		//),
		mcp.WithString("filters",
			mcp.Description("The filters to apply to the metrics. It is a stringified map[string]string[], e.g., '{\"service.name\": [\"serviceX\"]}' should return metrics for serviceX"),
		),
		mcp.WithString("excludeFilters",
			mcp.Description("The filters that should be excluded from the metrics. It is a stringified map[string]string[] e.g., '{\"service.name\": [\"serviceX\"]}' should return all metrics except for serviceX"),
		),
		mcp.WithString("splits",
			mcp.Description("JSON array of attributes to split the metrics by, e.g., '[\"service_name\"]' to split metrics by service"),
		),
		mcp.WithString("aggregation",
			mcp.Description("The operation to apply to the metrics (e.g., 'sum', 'avg')"),
		),
		mcp.WithString("functions",
			mcp.Description("JSON array of functions to apply to the metric, with each function having a name and parameters"),
		),
		mcp.WithBoolean("limitResults",
			mcp.Description("Flag to indicate if the results should be limited"),
		),
		mcp.WithNumber("bucketSize",
			mcp.Description("Size of each datapoint bucket in seconds"),
		),
	),
}

func main() {
	// Check if the appropriate environment variables are set
	if err := checkEnvVars(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

	s := server.NewMCPServer(
		"Metoro MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register all tools
	for _, tool := range tools {
		s.AddTool(tool, handlers[tool.Name])
	}

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
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
