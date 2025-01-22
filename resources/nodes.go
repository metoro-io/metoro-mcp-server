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

func NodesResourceHandler() (*mcpgolang.ResourceResponse, error) {
	now := time.Now()
	tenMinsAgo := now.Add(-10 * time.Minute)
	request := model.GetAllNodesRequest{
		StartTime:      tenMinsAgo.Unix(),
		EndTime:        now.Unix(),
		Filters:        map[string][]string{},
		ExcludeFilters: map[string][]string{},
		Splits:         []string{},
		Environments:   []string{},
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	response, err := utils.MakeMetoroAPIRequest("POST", "infrastructure/nodes", bytes.NewBuffer(jsonRequest), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}
	return mcpgolang.NewResourceResponse(
		mcpgolang.NewTextEmbeddedResource("api://nodes", string(response), "text/plain")), nil
}
