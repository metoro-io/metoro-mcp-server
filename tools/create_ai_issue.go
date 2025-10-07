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

type CreateAIIssueHandlerArgs struct {
	Title       string `json:"title" jsonschema:"required,description=Title of the AI issue"`
	Description string `json:"description" jsonschema:"required,description=Detailed description of the AI issue"`
}

func CreateAIIssueHandler(ctx context.Context, arguments CreateAIIssueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	request := model.CreateAIIssueRequest{
		Title:       arguments.Title,
		Description: arguments.Description,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	responseBody, err := utils.MakeMetoroAPIRequest("POST", "aiIssue", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to create AI issue: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
