package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func ListAIIssuesHandler(ctx context.Context) (*mcpgolang.ToolResponse, error) {
	responseBody, err := utils.MakeMetoroAPIRequest("GET", "aiIssues", nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to list AI issues: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
