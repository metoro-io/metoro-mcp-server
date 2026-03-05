package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

const (
	investigationCategoryDeploymentVerification = "deployment_verification"
	investigationCategoryAnomalyInvestigation   = "anomaly_investigation"
	investigationCategoryAlertInvestigation     = "alert_investigation"
)

type CreateInvestigationHandlerArgs struct {
	Title                   string           `json:"title" jsonschema:"required,description=Title of the investigation"`
	Category                string           `json:"category" jsonschema:"required,enum=deployment_verification,enum=anomaly_investigation,enum=alert_investigation,description=Category of investigation"`
	Verdict                 *string          `json:"verdict,omitempty" jsonschema:"description=Optional verdict for the investigation. Required when category is deployment_verification."`
	Summary                 string           `json:"summary" jsonschema:"description=Summary of the investigation - should be at most 3 sentences"`
	RecommendedActions      *[]string        `json:"recommendedActions,omitempty" jsonschema:"description=Optional recommended actions to take to remedy the issue. Should be concise - each item should be a single sentence."`
	ServiceName             *string          `json:"serviceName,omitempty" jsonschema:"description=Optional root cause service name to associate with this investigation."`
	Environment             *string          `json:"environment,omitempty" jsonschema:"description=Optional environment to associate with this investigation (e.g. production or staging)."`
	Namespace               *string          `json:"namespace,omitempty" jsonschema:"description=Optional Kubernetes namespace to associate with this investigation."`
	Markdown                string           `json:"markdown" jsonschema:"required,description=Markdown content of the investigation"`
	InProgress              *bool            `json:"inProgress" jsonschema:"description=Whether the investigation is in progress or not. Defaults to false"`
	TimeConfig              utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get the pods for. e.g. if you want the get the pods for the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	ChatHistoryUUID         *string          `json:"chatHistoryUuid,omitempty" jsonschema:"description=Optional chat history UUID to associate with this investigation"`
	IssueUUID               *string          `json:"issueUuid,omitempty" jsonschema:"description=Optional related AI issue UUID for this investigation"`
	AlertFireUUID           *string          `json:"alertFireUuid,omitempty" jsonschema:"description=Optional alert fire UUID to associate with this investigation"`
	AlertUUID               *string          `json:"alertUuid,omitempty" jsonschema:"description=Optional alert UUID to associate with this investigation"`
	DeploymentEventUUID     *string          `json:"deploymentEventUuid,omitempty" jsonschema:"description=Optional deployment event UUID to associate with this investigation for notification threading"`
	PotentialIssueEventUUID *string          `json:"potentialIssueEventUuid,omitempty" jsonschema:"description=Optional potential issue event UUID to associate with this investigation for notification threading"`
}

func validateInvestigationCategoryAndVerdict(category string, verdict *string) (*string, error) {
	switch category {
	case investigationCategoryDeploymentVerification, investigationCategoryAnomalyInvestigation, investigationCategoryAlertInvestigation:
	default:
		return nil, fmt.Errorf("invalid category: must be one of deployment_verification, anomaly_investigation, or alert_investigation")
	}

	if category == investigationCategoryDeploymentVerification {
		if verdict == nil || strings.TrimSpace(*verdict) == "" {
			return nil, fmt.Errorf("verdict is required when category is deployment_verification")
		}
	}

	if verdict == nil {
		return nil, nil
	}

	trimmedVerdict := strings.TrimSpace(*verdict)
	if trimmedVerdict == "" {
		return nil, nil
	}

	return &trimmedVerdict, nil
}

func CreateInvestigationHandler(ctx context.Context, arguments CreateInvestigationHandlerArgs) (*mcpgolang.ToolResponse, error) {
	trimmedVerdict, err := validateInvestigationCategoryAndVerdict(arguments.Category, arguments.Verdict)
	if err != nil {
		return nil, err
	}

	// Create the request body
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	truePtr := true
	reviewRequiredPtr := "ReviewRequired"
	start := time.Unix(startTime, 0)
	end := time.Unix(endTime, 0)
	tags := make(map[string]string)
	if arguments.ServiceName != nil {
		tags["service"] = *arguments.ServiceName
	}
	if arguments.Category == investigationCategoryDeploymentVerification {
		tags["verdict"] = *trimmedVerdict
	}
	request := model.CreateInvestigationRequest{
		Title:                   arguments.Title,
		Category:                arguments.Category,
		Verdict:                 trimmedVerdict,
		Summary:                 arguments.Summary,
		RecommendedActions:      arguments.RecommendedActions,
		Markdown:                arguments.Markdown,
		Tags:                    tags,
		IssueStartTime:          &start,
		IssueEndTime:            &end,
		ChatHistoryUUID:         arguments.ChatHistoryUUID,
		IsVisible:               &truePtr,
		InProgress:              arguments.InProgress,
		MetoroApprovalStatus:    &reviewRequiredPtr,
		IssueUUID:               arguments.IssueUUID,
		AlertFireUUID:           arguments.AlertFireUUID,
		AlertUUID:               arguments.AlertUUID,
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

	// Make the API request
	responseBody, err := utils.MakeMetoroAPIRequest("POST", "investigation", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to create investigation: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
