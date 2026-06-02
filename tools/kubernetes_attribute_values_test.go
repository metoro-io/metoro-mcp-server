package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TestGetAttributeValuesHandlerSupportsKubernetesResourceAliases(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"
	expectedStart := mustParseRFC3339Unix(t, start)
	expectedEnd := mustParseRFC3339Unix(t, end)

	var mu sync.Mutex
	var keysRequest *model.MultiMetricAttributeKeysRequest
	var valuesRequest *model.GetAttributeValuesRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/metrics/attributes":
			var req model.MultiMetricAttributeKeysRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode attribute keys request: %v", err)
			}
			mu.Lock()
			copied := req
			keysRequest = &copied
			mu.Unlock()
			_, _ = w.Write([]byte(`{"attributes":["Environment","Namespace","Kind","ResourceName","ServiceName","Reason","pod_phase"]}`))
		case "/api/v1/metrics/attribute/values":
			var req model.GetAttributeValuesRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode attribute values request: %v", err)
			}
			mu.Lock()
			copied := req
			valuesRequest = &copied
			mu.Unlock()
			_, _ = w.Write([]byte(`{"attribute":[{"value":"prod-azure","volume":1}]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := GetAttributeValuesHandler(context.Background(), GetAttributeValuesHandlerArgs{
		Type: model.MetricType("kubernetes_resources"),
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		Attribute: "environment",
		Filters: []model.Filter{
			{Key: "namespace", Values: []string{"metoro"}},
			{Key: "service.name", Values: []string{"/k8s/metoro/apiserver"}},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if keysRequest == nil {
		t.Fatalf("expected attribute keys request")
	}
	if keysRequest.Type != string(model.KubernetesResource) {
		t.Fatalf("expected canonical type %q, got %q", model.KubernetesResource, keysRequest.Type)
	}

	if valuesRequest == nil {
		t.Fatalf("expected attribute values request")
	}
	if valuesRequest.Type != model.KubernetesResource {
		t.Fatalf("expected values type %q, got %q", model.KubernetesResource, valuesRequest.Type)
	}
	if valuesRequest.Attribute != "Environment" {
		t.Fatalf("expected normalized attribute Environment, got %q", valuesRequest.Attribute)
	}
	if valuesRequest.Kubernetes == nil {
		t.Fatalf("expected kubernetes request body")
	}
	if valuesRequest.Kubernetes.StartTime != expectedStart {
		t.Fatalf("expected startTime %d, got %d", expectedStart, valuesRequest.Kubernetes.StartTime)
	}
	if valuesRequest.Kubernetes.EndTime != expectedEnd {
		t.Fatalf("expected endTime %d, got %d", expectedEnd, valuesRequest.Kubernetes.EndTime)
	}
	if got := strings.Join(valuesRequest.Kubernetes.Filters["Namespace"], ","); got != "metoro" {
		t.Fatalf("expected normalized Namespace filter, got %q", got)
	}
	if got := strings.Join(valuesRequest.Kubernetes.Filters["ServiceName"], ","); got != "/k8s/metoro/apiserver" {
		t.Fatalf("expected normalized ServiceName filter, got %q", got)
	}
}

func TestGetAttributeValuesHandlerRejectsInvalidKubernetesFilter(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"
	valuesCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/metrics/attributes":
			_, _ = w.Write([]byte(`{"attributes":["Environment","Namespace","Kind","ResourceName","ServiceName","Reason","pod_phase"]}`))
		case "/api/v1/metrics/attribute/values":
			valuesCalled = true
			_, _ = w.Write([]byte(`{"attribute":[]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := GetAttributeValuesHandler(context.Background(), GetAttributeValuesHandlerArgs{
		Type: model.KubernetesResource,
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		Attribute: "Environment",
		Filters: []model.Filter{
			{Key: "not_a_kubernetes_attribute", Values: []string{"x"}},
		},
	})
	if err == nil {
		t.Fatalf("expected invalid filter error")
	}
	if !strings.Contains(err.Error(), "invalid filter key") {
		t.Fatalf("expected invalid filter key error, got %v", err)
	}
	if valuesCalled {
		t.Fatalf("attribute values endpoint should not be called when validation fails")
	}
}

func TestGetAttributeKeysHandlerNormalizesKubernetesResourcesType(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	var captured *model.MultiMetricAttributeKeysRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/metrics/attributes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req model.MultiMetricAttributeKeysRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode attribute keys request: %v", err)
		}
		copied := req
		captured = &copied
		_, _ = w.Write([]byte(`{"attributes":["Environment"]}`))
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := GetAttributeKeysHandler(context.Background(), GetAttributeKeysHandlerArgs{
		Type: model.MetricType("kubernetes_resources"),
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if captured == nil {
		t.Fatalf("expected attribute keys request")
	}
	if captured.Type != string(model.KubernetesResource) {
		t.Fatalf("expected canonical type %q, got %q", model.KubernetesResource, captured.Type)
	}
}

func TestConvertTimeseriesToAPITimeseriesNormalizesKubernetesAliases(t *testing.T) {
	converted := convertTimeseriesToAPITimeseries([]model.SingleTimeseriesRequest{
		{
			Type: model.MetricType("kubernetes_resources"),
			Filters: []model.Filter{
				{Key: "namespace", Values: []string{"metoro"}},
				{Key: "service.name", Values: []string{"/k8s/metoro/apiserver"}},
			},
			ExcludeFilters: []model.Filter{
				{Key: "resourceName", Values: []string{"old-pod"}},
			},
			Splits:      []string{"service.name", "resourceName"},
			Aggregation: "count",
		},
	}, 100, 200)

	if len(converted) != 1 {
		t.Fatalf("expected one converted timeseries, got %d", len(converted))
	}
	metric := converted[0]
	if metric.Type != string(model.KubernetesResource) {
		t.Fatalf("expected canonical type %q, got %q", model.KubernetesResource, metric.Type)
	}
	if metric.KubernetesResource == nil {
		t.Fatalf("expected kubernetes resource request")
	}
	if got := strings.Join(metric.KubernetesResource.Filters["Namespace"], ","); got != "metoro" {
		t.Fatalf("expected Namespace filter, got %q", got)
	}
	if got := strings.Join(metric.KubernetesResource.Filters["ServiceName"], ","); got != "/k8s/metoro/apiserver" {
		t.Fatalf("expected ServiceName filter, got %q", got)
	}
	if got := strings.Join(metric.KubernetesResource.ExcludeFilters["ResourceName"], ","); got != "old-pod" {
		t.Fatalf("expected ResourceName exclude filter, got %q", got)
	}
	if strings.Join(metric.KubernetesResource.Splits, ",") != "ServiceName,ResourceName" {
		t.Fatalf("expected normalized splits, got %v", metric.KubernetesResource.Splits)
	}
}
