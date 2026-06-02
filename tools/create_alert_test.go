package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TestCreateAlertHandlerArgsUnmarshalsInvestigateOnFire(t *testing.T) {
	var args CreateAlertHandlerArgs
	if err := json.Unmarshal([]byte(`{"metoroql":"max(container_cpu_usage_seconds_total)","bucket_size":120,"investigate_on_fire":true}`), &args); err != nil {
		t.Fatalf("failed to unmarshal args: %v", err)
	}

	if args.MetoroQL != "max(container_cpu_usage_seconds_total)" {
		t.Fatalf("expected metoroql to unmarshal, got %q", args.MetoroQL)
	}
	if args.BucketSize != 120 {
		t.Fatalf("expected bucket_size to unmarshal to 120, got %d", args.BucketSize)
	}
	if !args.InvestigateOnFire {
		t.Fatal("expected investigate_on_fire to unmarshal to true")
	}
}

func TestCreateAlertPayloadUsesTopLevelInvestigateOnFire(t *testing.T) {
	alert := model.Alert{
		Metadata: model.MetadataObject{
			Name: "High 5XX Rate",
			Id:   "high-5xx-rate",
		},
		Timeseries: model.TimeseriesConfig{},
	}
	alert.SetInvestigateOnFire(true)

	request := model.CreateUpdateAlertRequest{Alert: alert}
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestJSON, &payload); err != nil {
		t.Fatalf("failed to unmarshal request JSON: %v", err)
	}
	alertPayload, ok := payload["alert"].(map[string]any)
	if !ok {
		t.Fatalf("expected alert payload object, got %T", payload["alert"])
	}

	if got := alertPayload["investigateOnFire"]; got != true {
		t.Fatalf("expected top-level alert.investigateOnFire=true, got %#v", got)
	}
	if metadataPayload, ok := alertPayload["metadata"].(map[string]any); ok {
		if _, exists := metadataPayload["investigateOnFire"]; exists {
			t.Fatal("did not expect investigateOnFire under alert.metadata")
		}
	}
	if timeseriesPayload, ok := alertPayload["timeseries"].(map[string]any); ok {
		if _, exists := timeseriesPayload["investigateOnFire"]; exists {
			t.Fatal("did not expect investigateOnFire under alert.timeseries")
		}
	}
}

func TestNewAlertOmitsInvestigateOnFireByDefault(t *testing.T) {
	alert := model.NewAlert(model.MetadataObject{
		Name: "High 5XX Rate",
		Id:   "high-5xx-rate",
	}, model.TimeseriesConfig{})

	request := model.CreateUpdateAlertRequest{Alert: *alert}
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestJSON, &payload); err != nil {
		t.Fatalf("failed to unmarshal request JSON: %v", err)
	}
	alertPayload, ok := payload["alert"].(map[string]any)
	if !ok {
		t.Fatalf("expected alert payload object, got %T", payload["alert"])
	}
	if _, exists := alertPayload["investigateOnFire"]; exists {
		t.Fatal("did not expect investigateOnFire to be marshaled by default")
	}
}

func TestCreateAlertHandlerCreatesMetoroQLAlertAfterValidation(t *testing.T) {
	query := `max(container_cpu_usage_seconds_total{namespace="default"}) by (service_name)`

	var fuzzyRequest model.FuzzyMetricsRequest
	var alertPayload map[string]any
	alertUpdateCalled := false

	withFakeMetoroAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/metoroql/convert/metoroqlToMetricSpecifier":
			var request metoroQLQueriesRequest
			decodeJSONBody(t, r, &request)
			if len(request.Queries) != 1 || request.Queries[0] != query {
				http.Error(w, fmt.Sprintf("unexpected queries: %#v", request.Queries), http.StatusBadRequest)
				return
			}
			writeJSON(w, `{
				"metricSpecifiers": [
					{
						"metricType": "metric",
						"metricName": "container_cpu_usage_seconds_total",
						"filters": {"namespace": ["default"]},
						"splits": ["service_name"],
						"aggregation": "max",
						"bucketSize": 60
					}
				],
				"formulas": [{"formula": "a"}]
			}`)
		case "/api/v1/fuzzyMetricsNames":
			decodeJSONBody(t, r, &fuzzyRequest)
			writeJSON(w, `{"metrics":["container_cpu_usage_seconds_total"]}`)
		case "/api/v1/metrics/attributes":
			writeJSON(w, `{"attributes":["namespace","service_name"]}`)
		case "/api/v1/alerts/update":
			alertUpdateCalled = true
			decodeJSONBody(t, r, &alertPayload)
			writeJSON(w, `{"id":"created-alert"}`)
		default:
			http.NotFound(w, r)
		}
	})

	_, err := CreateAlertHandler(context.Background(), CreateAlertHandlerArgs{
		AlertName:         "Container CPU",
		AlertDescription:  "High container CPU",
		MetoroQL:          query,
		Condition:         "GreaterThan",
		Threshold:         80,
		DatapointsToAlarm: 2,
		EvaluationWindow:  3,
		InvestigateOnFire: true,
	})
	if err != nil {
		t.Fatalf("expected create_alert to succeed, got %v", err)
	}

	if !alertUpdateCalled {
		t.Fatal("expected alerts/update to be called")
	}
	if fuzzyRequest.MetricFuzzyMatch != "container_cpu_usage_seconds_total" {
		t.Fatalf("expected targeted metric lookup, got %q", fuzzyRequest.MetricFuzzyMatch)
	}
	if fuzzyRequest.Discovery {
		t.Fatal("expected metric lookup discovery=false")
	}

	metoroQLTimeseries := alertPayload["alert"].(map[string]any)["timeseries"].(map[string]any)["expression"].(map[string]any)["metoroQLTimeseries"].(map[string]any)
	if got := metoroQLTimeseries["query"]; got != query {
		t.Fatalf("expected query %q, got %#v", query, got)
	}
	if got := metoroQLTimeseries["bucketSize"]; got != float64(60) {
		t.Fatalf("expected default bucketSize=60, got %#v", got)
	}
	if got := alertPayload["alert"].(map[string]any)["investigateOnFire"]; got != true {
		t.Fatalf("expected investigateOnFire=true, got %#v", got)
	}
}

func TestCreateAlertHandlerReturnsMetoroQLErrorBeforeCreatingAlert(t *testing.T) {
	alertUpdateCalled := false

	withFakeMetoroAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/metoroql/convert/metoroqlToMetricSpecifier":
			http.Error(w, `{"error":"Invalid MetoroQL query: parse error"}`, http.StatusBadRequest)
		case "/api/v1/alerts/update":
			alertUpdateCalled = true
			http.Error(w, "should not create alert", http.StatusInternalServerError)
		default:
			http.NotFound(w, r)
		}
	})

	_, err := CreateAlertHandler(context.Background(), CreateAlertHandlerArgs{
		AlertName:         "Broken query",
		AlertDescription:  "Broken query",
		MetoroQL:          "max(",
		Condition:         "GreaterThan",
		Threshold:         1,
		DatapointsToAlarm: 1,
		EvaluationWindow:  1,
	})
	if err == nil || !strings.Contains(err.Error(), "invalid MetoroQL query") {
		t.Fatalf("expected invalid MetoroQL query error, got %v", err)
	}
	if alertUpdateCalled {
		t.Fatal("did not expect alerts/update to be called")
	}
}

func TestCreateAlertHandlerReturnsTimeseriesValidationErrorBeforeCreatingAlert(t *testing.T) {
	alertUpdateCalled := false

	withFakeMetoroAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/metoroql/convert/metoroqlToMetricSpecifier":
			writeJSON(w, `{
				"metricSpecifiers": [
					{
						"metricType": "metric",
						"metricName": "container_cpu_usage_seconds_total",
						"filters": {"namespace": ["default"]},
						"aggregation": "max",
						"bucketSize": 60
					}
				]
			}`)
		case "/api/v1/fuzzyMetricsNames":
			writeJSON(w, `{"metrics":["container_cpu_usage_seconds_total"]}`)
		case "/api/v1/metrics/attributes":
			writeJSON(w, `{"attributes":["pod_name"]}`)
		case "/api/v1/alerts/update":
			alertUpdateCalled = true
			http.Error(w, "should not create alert", http.StatusInternalServerError)
		default:
			http.NotFound(w, r)
		}
	})

	_, err := CreateAlertHandler(context.Background(), CreateAlertHandlerArgs{
		AlertName:         "Bad attribute",
		AlertDescription:  "Bad attribute",
		MetoroQL:          `max(container_cpu_usage_seconds_total{namespace="default"})`,
		Condition:         "GreaterThan",
		Threshold:         1,
		DatapointsToAlarm: 1,
		EvaluationWindow:  1,
	})
	if err == nil || !strings.Contains(err.Error(), "invalid filter key: namespace") {
		t.Fatalf("expected invalid filter key error, got %v", err)
	}
	if alertUpdateCalled {
		t.Fatal("did not expect alerts/update to be called")
	}
}

func TestCreateAlertFromMetoroQLRejectsNegativeBucketSize(t *testing.T) {
	_, err := createAlertFromMetoroQL(context.Background(), "test", "test", "max(metric_name)", -1, "GreaterThan", 1, 1, 1, false)
	if err == nil || !strings.Contains(err.Error(), "bucket_size must be positive") {
		t.Fatalf("expected bucket_size error, got %v", err)
	}
}

func withFakeMetoroAPI(t *testing.T, handler http.HandlerFunc) {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	t.Setenv(utils.METORO_API_URL_ENV_VAR, server.URL)
	t.Setenv(utils.METORO_AUTH_TOKEN_ENV_VAR, "test-token")
}

func decodeJSONBody(t *testing.T, r *http.Request, target any) {
	t.Helper()

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode %s request body: %v", r.URL.Path, err)
	}
}

func writeJSON(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}
