package tools

import (
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type ListAIIssueEventsHandlerArgs struct {
	IssueUUID string `json:"issueUuid" jsonschema:"required,description=UUID of the AI issue whose events should be listed"`
}

func ListAIIssueEventsHandler(ctx context.Context, arguments ListAIIssueEventsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	endpoint := fmt.Sprintf("aiIssue/events?issueUuid=%s", arguments.IssueUUID)
	responseBody, err := utils.MakeMetoroAPIRequest("GET", endpoint, nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to list AI issue events: %w", err)
	}

	var resp model.ListAIIssueEventsResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse AI issue events response: %w", err)
	}

	serialized, err := json.Marshal(resp.Events)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AI issue events: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(serialized))), nil
}
