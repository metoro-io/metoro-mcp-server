package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type CreateDashboardHandlerArgs struct {
	DashboardName string            `json:"dashboard_name" jsonschema:"required,description=The name of the dashboard to create"`
	GroupWidget   model.GroupWidget `json:"group_widget" jsonschema:"required,description=The group widget this dashboard will have. This is the top level widget of the dashboard that will contain all other widgets. A widget can be either a group widget or a MetricChartWidget"`
}

func CreateDashboardHandler(ctx context.Context, arguments CreateDashboardHandlerArgs) (*mcpgolang.ToolResponse, error) {
	dashboardJson, err := json.Marshal(arguments.GroupWidget)
	if err != nil {
		return nil, fmt.Errorf("error marshaling dashboard properties: %v", err)
	}

	newDashboardRequest := model.SetDashboardRequest{
		Name:             arguments.DashboardName,
		Id:               uuid.NewString(),
		DashboardJson:    string(dashboardJson),
		DefaultTimeRange: "1h",
	}

	resp, err := setDashboardMetoroCall(ctx, newDashboardRequest)
	if err != nil {
		return nil, fmt.Errorf("error setting dashboard: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func setDashboardMetoroCall(ctx context.Context, request model.SetDashboardRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling dashboard request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "dashboard", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
