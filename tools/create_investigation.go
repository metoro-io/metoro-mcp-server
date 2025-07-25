package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type CreateInvestigationHandlerArgs struct {
	Title                   string           `json:"title" jsonschema:"required,description=Title of the investigation"`
	Markdown                string           `json:"markdown" jsonschema:"required,description=Markdown content of the investigation"`
	TimeConfig              utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get the pods for. e.g. if you want the get the pods for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	ChatHistoryUUID         *string          `json:"chatHistoryUuid,omitempty" jsonschema:"description=Optional chat history UUID to associate with this investigation"`
	ParentInvestigationUUID *string          `json:"parentInvestigationUuid,omitempty" jsonschema:"description=Optional parent investigation UUID to associate with this investigation. Set this if this is a recurrence of an existing investigation"`
}

func CreateInvestigationHandler(ctx context.Context, arguments CreateInvestigationHandlerArgs) (*mcpgolang.ToolResponse, error) {
	// Create the request body
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	truePtr := true
	reviewRequiredPtr := "ReviewRequired"
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)
	request := model.CreateInvestigationRequest{
		Title:                   arguments.Title,
		Markdown:                arguments.Markdown,
		IssueStartTime:          &start,
		IssueEndTime:            &end,
		ChatHistoryUUID:         arguments.ChatHistoryUUID,
		IsVisible:               &truePtr,
		MetoroApprovalStatus:    &reviewRequiredPtr,
		ParentInvestigationUUID: arguments.ParentInvestigationUUID,
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
