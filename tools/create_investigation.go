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
	Summary                 string           `json:"summary" jsonschema:"description=Summary of the investigation - should be at most 3 sentences"`
	RecommendedActions      *[]string        `json:"recommendedActions,omitempty" jsonschema:"description=Optional recommended actions to take to remedy the issue. Should be concise - each item should be a single sentence."`
	ServiceName             *string          `json:"serviceName,omitempty" jsonschema:"description=Optional service name to associate with this investigation. If you are investigating an issue in a specific service you should set this. If set the investigation will be associated with the service. "`
	Markdown                string           `json:"markdown" jsonschema:"required,description=Markdown content of the investigation"`
	InProgress              *bool            `json:"inProgress" jsonschema:"description=Whether the investigation is in progress or not. Defaults to false"`
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

	falsePtr := false
	reviewRequiredPtr := "ReviewRequired"
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)
	request := model.CreateInvestigationRequest{
		Title:                   arguments.Title,
		Summary:                 arguments.Summary,
		RecommendedActions:      arguments.RecommendedActions,
		Markdown:                arguments.Markdown,
		ServiceName:             arguments.ServiceName,
		IssueStartTime:          &start,
		IssueEndTime:            &end,
		ChatHistoryUUID:         arguments.ChatHistoryUUID,
		IsVisible:               &falsePtr,
		InProgress:              arguments.InProgress,
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
