package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model/publicapi"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type CreateDashboardHandlerArgs struct {
	Dashboard DashboardPayload `json:"dashboard" jsonschema:"required,description=Dashboard object matching Metoro's public Dashboard OpenAPI schema. Required fields are dashboard.metadata.id and dashboard.content. dashboard.content is a group widget; nested widgets may use type group, markdown, chart, stat, gauge, trace, or log."`
}

type DashboardPayload json.RawMessage

type dashboardPayloadSchema struct {
	Metadata map[string]any `json:"metadata" jsonschema:"required,description=Dashboard metadata. metadata.id is required. metadata.title and metadata.folderPath are optional."`
	Content  map[string]any `json:"content" jsonschema:"required,description=Top-level group widget from Metoro's public Dashboard OpenAPI schema."`
	Settings map[string]any `json:"settings,omitempty" jsonschema:"description=Optional dashboard settings from Metoro's public Dashboard OpenAPI schema."`
}

func (DashboardPayload) JSONSchemaAlias() any {
	return dashboardPayloadSchema{}
}

func (p DashboardPayload) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return []byte("null"), nil
	}
	if !json.Valid(p) {
		return nil, fmt.Errorf("invalid dashboard JSON")
	}
	return append([]byte(nil), p...), nil
}

func (p *DashboardPayload) UnmarshalJSON(data []byte) error {
	if !json.Valid(data) {
		return fmt.Errorf("invalid dashboard JSON")
	}
	*p = append((*p)[0:0], data...)
	return nil
}

func (a *CreateDashboardHandlerArgs) UnmarshalJSON(data []byte) error {
	var raw struct {
		Dashboard DashboardPayload `json:"dashboard"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&raw); err != nil {
		return err
	}
	if len(raw.Dashboard) == 0 || strings.TrimSpace(string(raw.Dashboard)) == "null" {
		return fmt.Errorf("dashboard is required")
	}

	var request publicapi.CreateUpdateDashboardRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}

	a.Dashboard = raw.Dashboard
	return nil
}

func CreateDashboardHandler(ctx context.Context, arguments CreateDashboardHandlerArgs) (*mcpgolang.ToolResponse, error) {
	resp, err := setDashboardMetoroCall(ctx, arguments)
	if err != nil {
		return nil, fmt.Errorf("error creating or updating dashboard: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func setDashboardMetoroCall(ctx context.Context, request CreateDashboardHandlerArgs) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling dashboard request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "dashboards/update", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
