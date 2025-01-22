package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func MetricsResourceHandler() (*mcpgolang.ResourceResponse, error) {
	now := time.Now()
	twoHoursAgo := now.Add(-2 * time.Hour)
	request := model.FuzzyMetricsRequest{
		StartTime:        twoHoursAgo.Unix(),
		EndTime:          now.Unix(),
		MetricFuzzyMatch: "",         // This will return all the metric names.
		Environments:     []string{}, // All environments
	}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	resp, err := utils.MakeMetoroAPIRequest("POST", "fuzzyMetricsNames", bytes.NewBuffer(jsonData), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}
	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://metrics", string(resp), "text/plain")), nil
}
