package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetMetricNamesHandlerArgs struct {
	Environments []string `json:"environments" jsonschema:"description=Environments to get metrics names from. If empty all environments will be used."`
}

func GetMetricNamesHandler(ctx context.Context, arguments GetMetricNamesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	hourAgo := now.Add(-1 * time.Hour)
	request := model.FuzzyMetricsRequest{
		StartTime:        hourAgo.Unix(),
		EndTime:          now.Unix(),
		MetricFuzzyMatch: "", // This will return all the metric names.
		Environments:     arguments.Environments,
	}
	response, err := getMetricNamesMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}

func getMetricNamesMetoroCall(ctx context.Context, request model.FuzzyMetricsRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return utils.MakeMetoroAPIRequest("POST", "fuzzyMetricsNames", bytes.NewBuffer(jsonData), utils.GetAPIRequirementsFromRequest(ctx))
}
