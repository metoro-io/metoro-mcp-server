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

type CreateInvestigationHandlerArgs struct {
	Title           string  `json:"title" jsonschema:"required,description=Title of the investigation"`
	Markdown        string  `json:"markdown" jsonschema:"required,description=Markdown content of the investigation"`
	IssueStartTime  *int64  `json:"issueStartTime,omitempty" jsonschema:"description=Optional start time of the issue in seconds since epoch"`
	IssueEndTime    *int64  `json:"issueEndTime,omitempty" jsonschema:"description=Optional end time of the issue in seconds since epoch"`
	ChatHistoryUUID *string `json:"chatHistoryUuid,omitempty" jsonschema:"description=Optional chat history UUID to associate with this investigation"`
}

func CreateInvestigationHandler(ctx context.Context, arguments CreateInvestigationHandlerArgs) (*mcpgolang.ToolResponse, error) {
	// Create the request body
	falsePtr := false
	reviewRequiredPtr := "ReviewRequired"
	request := model.CreateInvestigationRequest{
		Title:                arguments.Title,
		Markdown:             arguments.Markdown,
		IssueStartTime:       arguments.IssueStartTime,
		IssueEndTime:         arguments.IssueEndTime,
		ChatHistoryUUID:      arguments.ChatHistoryUUID,
		IsVisible:            &falsePtr,
		MetoroApprovalStatus: &reviewRequiredPtr,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	responseBody, err := utils.MakeMetoroAPIRequest("POST", "investigation", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to create investigation: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
