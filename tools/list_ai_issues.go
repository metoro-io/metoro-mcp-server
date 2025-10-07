package tools

import (
	"context"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type ListAIIssueHandlerArgs struct {
	OpenOnly *bool `json:"openOnly,omitempty" jsonschema:"description=Set to true to list only open issues (default true), false to list all issues"`
}

func ListAIIssuesHandler(ctx context.Context, arguments ListAIIssueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	openOnly := true
	if arguments.OpenOnly != nil {
		openOnly = *arguments.OpenOnly
	}

	endpoint := "aiIssues"
	if openOnly {
		endpoint += "?openOnly=true"
	}

	responseBody, err := utils.MakeMetoroAPIRequest("GET", endpoint, nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to list AI issues: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
