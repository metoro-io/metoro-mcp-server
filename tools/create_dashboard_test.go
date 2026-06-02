package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/invopop/jsonschema"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model/publicapi"
)

func TestCreateDashboardPayloadUsesPublicAPIShape(t *testing.T) {
	requestJSON, err := json.Marshal(sampleCreateDashboardArgs(t))
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestJSON, &payload); err != nil {
		t.Fatalf("failed to unmarshal request JSON: %v", err)
	}

	if _, exists := payload["dashboard"]; !exists {
		t.Fatal("expected public dashboard payload to contain top-level dashboard")
	}
	if _, exists := payload["name"]; exists {
		t.Fatal("did not expect legacy name field")
	}
	if _, exists := payload["dashboardJson"]; exists {
		t.Fatal("did not expect legacy dashboardJson field")
	}

	dashboard := payload["dashboard"].(map[string]any)
	metadata := dashboard["metadata"].(map[string]any)
	if got := metadata["id"]; got != "service-overview" {
		t.Fatalf("expected metadata.id service-overview, got %#v", got)
	}
	content := dashboard["content"].(map[string]any)
	widgets := content["widgets"].([]any)
	if len(widgets) != 1 {
		t.Fatalf("expected one widget, got %d", len(widgets))
	}
}

func TestSetDashboardMetoroCallUsesPublicUpdateEndpoint(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotAuth string
	var gotPayload map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if err := json.Unmarshal(body, &gotPayload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"service-overview"}`))
	}))
	defer server.Close()

	t.Setenv("METORO_API_URL", server.URL)
	t.Setenv("METORO_AUTH_TOKEN", "test-token")

	response, err := setDashboardMetoroCall(context.Background(), sampleCreateDashboardArgs(t))
	if err != nil {
		t.Fatalf("setDashboardMetoroCall returned error: %v", err)
	}

	if string(response) != `{"id":"service-overview"}` {
		t.Fatalf("unexpected response body: %s", string(response))
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/api/v1/dashboards/update" {
		t.Fatalf("expected /api/v1/dashboards/update, got %s", gotPath)
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("expected bearer auth header, got %q", gotAuth)
	}
	if _, exists := gotPayload["dashboard"]; !exists {
		t.Fatalf("expected request payload to contain dashboard, got %#v", gotPayload)
	}
}

func TestCreateDashboardToolRegistrationDoesNotRecurseThroughDashboardWidgets(t *testing.T) {
	server := mcpgolang.NewServer(nil)
	if err := server.RegisterTool("create_dashboard", "Create dashboard", CreateDashboardHandler); err != nil {
		t.Fatalf("failed to register create_dashboard: %v", err)
	}
}

func TestCreateDashboardInputSchemaKeepsDashboardAsObject(t *testing.T) {
	reflector := jsonschema.Reflector{
		Anonymous:                  true,
		AllowAdditionalProperties:  true,
		RequiredFromJSONSchemaTags: true,
		DoNotReference:             true,
		ExpandedStruct:             true,
	}
	schema := reflector.ReflectFromType(reflect.TypeOf(CreateDashboardHandlerArgs{}))
	if schema.Properties == nil {
		t.Fatalf("expected schema properties, got %#v", schema)
	}
	dashboard, ok := schema.Properties.Get("dashboard")
	if !ok {
		t.Fatal("expected dashboard property in schema")
	}

	schemaJSON, err := json.Marshal(dashboard)
	if err != nil {
		t.Fatalf("failed to marshal dashboard schema: %v", err)
	}
	schemaString := string(schemaJSON)

	if !strings.Contains(schemaString, `"type":"object"`) {
		t.Fatalf("expected dashboard schema to be an object, got %s", schemaString)
	}
	if !strings.Contains(schemaString, `"metadata"`) {
		t.Fatalf("expected dashboard schema to include metadata, got %s", schemaString)
	}
	if !strings.Contains(schemaString, `"content"`) {
		t.Fatalf("expected dashboard schema to include content, got %s", schemaString)
	}
	if strings.Contains(schemaString, `"type":"array"`) {
		t.Fatalf("dashboard schema should not expose raw JSON bytes as an array: %s", schemaString)
	}
}

func TestCreateDashboardHandlerArgsRejectsLegacyPayload(t *testing.T) {
	var args CreateDashboardHandlerArgs
	err := json.Unmarshal([]byte(`{"dashboard_name":"Legacy","group_widget":{"widgetType":"Group","children":[]}}`), &args)
	if err == nil {
		t.Fatal("expected legacy dashboard payload to be rejected")
	}
	if !strings.Contains(err.Error(), "dashboard") {
		t.Fatalf("expected missing dashboard error, got: %v", err)
	}
}

func TestRepresentativeChartWidgetPayloadMarshalsPublicShape(t *testing.T) {
	requestJSON, err := json.Marshal(sampleCreateDashboardArgs(t))
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestJSON, &payload); err != nil {
		t.Fatalf("failed to unmarshal request JSON: %v", err)
	}

	dashboard := payload["dashboard"].(map[string]any)
	content := dashboard["content"].(map[string]any)
	widget := content["widgets"].([]any)[0].(map[string]any)
	if got := widget["type"]; got != "chart" {
		t.Fatalf("expected chart widget type, got %#v", got)
	}
	if _, exists := widget["chart"]; !exists {
		t.Fatalf("expected chart widget payload, got %#v", widget)
	}
	if _, exists := widget["widgetType"]; exists {
		t.Fatal("did not expect legacy widgetType field")
	}

	chart := widget["chart"].(map[string]any)
	expression := chart["expression"].(map[string]any)
	queries := expression["metoroQLQueries"].([]any)
	query := queries[0].(map[string]any)
	if got := query["query"]; got != `avg(container_cpu_usage_seconds_total)` {
		t.Fatalf("expected MetoroQL query, got %#v", got)
	}
}

func sampleCreateDashboardArgs(t *testing.T) CreateDashboardHandlerArgs {
	t.Helper()

	requestJSON, err := json.Marshal(sampleCreateDashboardRequest())
	if err != nil {
		t.Fatalf("failed to marshal public dashboard request: %v", err)
	}

	var args CreateDashboardHandlerArgs
	if err := json.Unmarshal(requestJSON, &args); err != nil {
		t.Fatalf("failed to unmarshal dashboard args: %v", err)
	}
	return args
}

func sampleCreateDashboardRequest() publicapi.CreateUpdateDashboardRequest {
	title := "Service Overview"
	folderPath := "/dashboards/default/"
	defaultTimeRange := "1h"
	chartTitle := "CPU Usage"
	chartType := "line"
	bucketSize := int64(60)
	label := "CPU"

	return publicapi.CreateUpdateDashboardRequest{
		Dashboard: publicapi.Dashboard{
			Metadata: publicapi.DashboardMetadata{
				Id:         "service-overview",
				Title:      &title,
				FolderPath: &folderPath,
			},
			Content: publicapi.GroupWidget{
				Widgets: []publicapi.Widget{
					{
						Type: "chart",
						Position: &publicapi.Position{
							Type: "absolute",
							Absolute: &publicapi.PositionAbsolute{
								X: 0,
								Y: 0,
								W: 6,
								H: 3,
							},
						},
						Chart: &publicapi.ChartWidget{
							Title:     &chartTitle,
							ChartType: &chartType,
							Expression: &publicapi.ChartWidgetExpression{
								MetoroQLQueries: []publicapi.MetoroQLQuery{
									{
										Query:      `avg(container_cpu_usage_seconds_total)`,
										BucketSize: &bucketSize,
										DisplaySettings: &publicapi.ChartDisplaySettings{
											Label: &label,
										},
									},
								},
							},
						},
					},
				},
			},
			Settings: &publicapi.DashboardSettings{
				DefaultTimeRange: &defaultTimeRange,
			},
		},
	}
}
