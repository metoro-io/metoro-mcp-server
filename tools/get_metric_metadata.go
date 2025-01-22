package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetMetricMetadataHandlerArgs struct {
	Name string `json:"name" jsonschema:"required,description=The name of the metric to get metadata for"`
}

func GetMetricMetadata(ctx context.Context, arguments GetMetricMetadataHandlerArgs) (*mcpgolang.ToolResponse, error) {
	response, err := getMetricMetadataMetoroCall(ctx, arguments.Name)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro get metric metadata api: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}

func getMetricMetadataMetoroCall(ctx context.Context, metricName string) ([]byte, error) {
	return utils.MakeMetoroAPIRequest("GET", "metric/metadata?name="+metricName, nil, utils.GetAPIRequirementsFromRequest(ctx))
}
