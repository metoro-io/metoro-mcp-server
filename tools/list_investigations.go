package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type ListInvestigationsHandlerArgs struct {
	Limit           int               `json:"limit,omitempty" jsonschema:"description=Maximum number of investigations to return (default 20 max 100)"`
	Offset          int               `json:"offset,omitempty" jsonschema:"description=Number of investigations to skip for pagination"`
	Tags            map[string]string `json:"tags,omitempty" jsonschema:"description=Filter investigations by tags"`
	IncludeResolved bool              `json:"includeResolved,omitempty" jsonschema:"description=Include resolved investigations in the results"`
}

func ListInvestigationsHandler(ctx context.Context, arguments ListInvestigationsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	// Create the request body
	request := struct {
		Limit             int               `json:"limit,omitempty"`
		Offset            int               `json:"offset,omitempty"`
		Tags              map[string]string `json:"tags,omitempty"`
		IncludeResolved   bool              `json:"includeResolved,omitempty"`
		ExcludeInProgress bool              `json:"excludeInProgress,omitempty"`
	}{
		Limit:             arguments.Limit,
		Offset:            arguments.Offset,
		Tags:              arguments.Tags,
		IncludeResolved:   arguments.IncludeResolved,
		ExcludeInProgress: true, // Always exclude in-progress investigations as the AI only wants to see the past investigations.
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	responseBody, err := utils.MakeMetoroAPIRequest("POST", "investigations/list", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to list investigations: %w", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(responseBody))), nil
}
