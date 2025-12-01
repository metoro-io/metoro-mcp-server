package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type ReportDeploymentVerdictHandlerArgs struct {
	DeploymentEventUUID string `json:"deployment_event_uuid" jsonschema:"required,description=The UUID of the deployment event being evaluated. This should be provided in the context of the deployment investigation."`
	ServiceName         string `json:"service_name" jsonschema:"required,description=Name of the service that was deployed"`
	Environment         string `json:"environment" jsonschema:"required,description=Environment where deployment occurred"`
	Namespace           string `json:"namespace" jsonschema:"required,description=Kubernetes namespace"`
	Verdict             string `json:"verdict" jsonschema:"required,enum=healthy,enum=degraded,enum=failed,description=The deployment health verdict - healthy means no issues detected - degraded means minor issues - failed means critical issues requiring immediate attention"`
	Summary             string `json:"summary" jsonschema:"required,description=Brief 1-2 sentence summary of the verdict"`
	Reason              string `json:"reason" jsonschema:"required,description=Detailed explanation of why this verdict was reached - including specific metrics/errors found"`
}

type CreateDeploymentVerdictRequest struct {
	DeploymentEventUUID string `json:"deployment_event_uuid"`
	ServiceName         string `json:"service_name"`
	Environment         string `json:"environment"`
	Namespace           string `json:"namespace"`
	Verdict             string `json:"verdict"`
	Summary             string `json:"summary"`
	Reason              string `json:"reason"`
}

func ReportDeploymentVerdictHandler(ctx context.Context, arguments ReportDeploymentVerdictHandlerArgs) (*mcpgolang.ToolResponse, error) {
	// Validate verdict
	if arguments.Verdict != "healthy" && arguments.Verdict != "degraded" && arguments.Verdict != "failed" {
		return nil, fmt.Errorf("invalid verdict: must be one of healthy, degraded, or failed")
	}

	// Create the request body
	request := CreateDeploymentVerdictRequest{
		DeploymentEventUUID: arguments.DeploymentEventUUID,
		ServiceName:         arguments.ServiceName,
		Environment:         arguments.Environment,
		Namespace:           arguments.Namespace,
		Verdict:             arguments.Verdict,
		Summary:             arguments.Summary,
		Reason:              arguments.Reason,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	responseBody, err := utils.MakeMetoroAPIRequest("POST", "deploymentVerdict", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to report deployment verdict: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
