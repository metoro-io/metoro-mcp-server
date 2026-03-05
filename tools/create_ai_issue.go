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
	Title        string   `json:"title" jsonschema:"required,description=Title of the AI issue"`
	Description  string   `json:"description" jsonschema:"required,description=Detailed description of the AI issue"`
	Summary      string   `json:"summary" jsonschema:"required,description=One sentence summary of the AI issue"`
	Environments []string `json:"environments,omitempty" jsonschema:"description=Optional list of environments related to this issue"`
	Services     []string `json:"services,omitempty" jsonschema:"description=Optional list of services related to this issue"`
	Priority     *string  `json:"priority,omitempty" jsonschema:"enum=P1,enum=P2,enum=P3,description=Optional issue priority"`
	Category     *string  `json:"category,omitempty" jsonschema:"enum=application,enum=infrastructure,description=Optional issue category"`
}

func CreateAIIssueHandler(ctx context.Context, arguments CreateAIIssueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	if arguments.Priority != nil && *arguments.Priority != "P1" && *arguments.Priority != "P2" && *arguments.Priority != "P3" {
		return nil, fmt.Errorf("invalid priority: must be one of P1, P2, or P3")
	}
	if arguments.Category != nil && *arguments.Category != "application" && *arguments.Category != "infrastructure" {
		return nil, fmt.Errorf("invalid category: must be one of application or infrastructure")
	}

	request := model.CreateAIIssueRequest{
		Title:        arguments.Title,
		Description:  arguments.Description,
		Summary:      arguments.Summary,
		Environments: arguments.Environments,
		Services:     arguments.Services,
		Priority:     arguments.Priority,
		Category:     arguments.Category,
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
