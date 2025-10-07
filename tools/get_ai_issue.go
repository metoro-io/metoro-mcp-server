package tools

import (
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetAIIssueHandlerArgs struct {
	IssueUUID string `json:"issueUuid" jsonschema:"required,description=UUID of the AI issue to retrieve"`
}

func GetAIIssueHandler(ctx context.Context, arguments GetAIIssueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	endpoint := fmt.Sprintf("aiIssue?uuid=%s", arguments.IssueUUID)
	responseBody, err := utils.MakeMetoroAPIRequest("GET", endpoint, nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch AI issue: %w", err)
	}

	var issueResponse model.GetAIIssueResponse
	if err := json.Unmarshal(responseBody, &issueResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI issue response: %w", err)
	}

	serialized, err := json.Marshal(issueResponse.Issue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AI issue: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(serialized))), nil
}
