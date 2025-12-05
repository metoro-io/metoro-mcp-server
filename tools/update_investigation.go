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

type UpdateInvestigationHandlerArgs struct {
	InvestigationUUID   string           `json:"investigationUuid" jsonschema:"required,description=UUID of the investigation to update"`
	Title               string           `json:"title" jsonschema:"required,description=Title of the investigation"`
	Summary             string           `json:"summary" jsonschema:"description=Summary of the investigation - should be at most 3 sentences"`
	RecommendedActions  *[]string        `json:"recommendedActions,omitempty" jsonschema:"description=Optional recommended actions to take to remedy the issue. Should be concise - each item should be a single sentence."`
	ServiceName         *string          `json:"serviceName,omitempty" jsonschema:"description=Optional root cause service name to associate with this investigation."`
	Environment         *string          `json:"environment,omitempty" jsonschema:"description=Optional environment to associate with this investigation (e.g. production or staging)."`
	Namespace           *string          `json:"namespace,omitempty" jsonschema:"description=Optional Kubernetes namespace to associate with this investigation."`
	Markdown            string           `json:"markdown" jsonschema:"required,description=Markdown content of the investigation"`
	InProgress          *bool            `json:"inProgress" jsonschema:"description=Whether the investigation is in progress or not. Defaults to false"`
	TimeConfig          utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get the pods for. e.g. if you want the get the pods for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	ChatHistoryUUID     *string          `json:"chatHistoryUuid,omitempty" jsonschema:"description=Optional chat history UUID to associate with this investigation"`
	IssueUUID           *string          `json:"issueUuid,omitempty" jsonschema:"description=Optional related AI issue UUID for this investigation"`
	DeploymentEventUUID     *string          `json:"deploymentEventUuid,omitempty" jsonschema:"description=Optional deployment event UUID to associate with this investigation for notification threading"`
	PotentialIssueEventUUID *string          `json:"potentialIssueEventUuid,omitempty" jsonschema:"description=Optional potential issue event UUID to associate with this investigation for notification threading"`
}

func UpdateInvestigationHandler(ctx context.Context, arguments UpdateInvestigationHandlerArgs) (*mcpgolang.ToolResponse, error) {
	// Create the request body
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	falsePtr := false
	reviewRequiredPtr := "ReviewRequired"
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)

	tags := make(map[string]string)
	if arguments.ServiceName != nil {
		tags["service"] = *arguments.ServiceName
	}

	title := arguments.Title
	summary := arguments.Summary
	markdown := arguments.Markdown
	tagsPtr := tags

	request := model.UpdateInvestigationRequest{
		Title:                &title,
		Summary:              &summary,
		Markdown:             &markdown,
		Tags:                 &tagsPtr,
		IssueStartTime:       &start,
		IssueEndTime:         &end,
		ChatHistoryUUID:      arguments.ChatHistoryUUID,
		IsVisible:            &falsePtr,
		InProgress:           arguments.InProgress,
		MetoroApprovalStatus: &reviewRequiredPtr,
		IssueUUID:               arguments.IssueUUID,
		RecommendedActions:      arguments.RecommendedActions,
		DeploymentEventUUID:     arguments.DeploymentEventUUID,
		PotentialIssueEventUUID: arguments.PotentialIssueEventUUID,
		Environment:             arguments.Environment,
		Namespace:               arguments.Namespace,
		ServiceName:             arguments.ServiceName,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request - using PUT method for update
	endpoint := fmt.Sprintf("investigation?uuid=%s", arguments.InvestigationUUID)
	responseBody, err := utils.MakeMetoroAPIRequest("PUT", endpoint, bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to update investigation: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
