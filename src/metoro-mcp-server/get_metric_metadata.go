package main

import (
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
)

type GetMetricMetadataHandlerArgs struct {
	Name string `json:"name" jsonschema:"required,description=The name of the metric to get metadata for"`
}

func getMetricMetadata(arguments GetMetricMetadataHandlerArgs) (*mcpgolang.ToolResponse, error) {
	response, err := getMetricMetadataMetoroCall(arguments.Name)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro get metric metadata api: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}

func getMetricMetadataMetoroCall(metricName string) ([]byte, error) {
	return MakeMetoroAPIRequest("GET", "metric/metadata?name="+metricName, nil)
}
