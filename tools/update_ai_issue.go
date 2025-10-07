package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type UpdateAIIssueHandlerArgs struct {
	IssueUUID   string  `json:"issueUuid" jsonschema:"required,description=UUID of the AI issue to update"`
	Title       *string `json:"title,omitempty" jsonschema:"description=Optional new title for the AI issue"`
	Description *string `json:"description,omitempty" jsonschema:"description=Optional new description for the AI issue"`
	Open        *bool   `json:"open,omitempty" jsonschema:"description=Optional flag to set whether the AI issue is open (true) or resolved (false)"`
}

func UpdateAIIssueHandler(ctx context.Context, arguments UpdateAIIssueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	if arguments.Title == nil && arguments.Description == nil && arguments.Open == nil {
		return nil, fmt.Errorf("at least one of title, description, or open must be provided to update an AI issue")
	}

	request := model.UpdateAIIssueRequest{
		Title:       arguments.Title,
		Description: arguments.Description,
		Open:        arguments.Open,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("aiIssue?uuid=%s", arguments.IssueUUID)
	responseBody, err := utils.MakeMetoroAPIRequest("PUT", endpoint, bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to update AI issue: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
