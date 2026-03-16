package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TestGetMultiMetricHandlerAcceptsLabelsWithoutChangingAPIRequest(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	var captured model.GetMultiMetricRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/metrics/attributes":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"attributes":[]}`))
		case "/api/v1/metrics":
			if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
				t.Fatalf("failed to decode metrics request: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"metrics":[]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := GetMultiMetricHandler(context.Background(), GetMultiMetricHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		Timeseries: []model.SingleTimeseriesRequest{
			{
				Type:              model.Trace,
				Aggregation:       model.AggregationCount,
				BucketSize:        60,
				FormulaIdentifier: "a",
				Label:             "Request volume",
			},
		},
		Formulas: []model.Formula{
			{
				Formula: "a",
				Label:   "Displayed formula label",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(captured.Metrics) != 1 {
		t.Fatalf("expected 1 metric request, got %d", len(captured.Metrics))
	}
	if captured.Metrics[0].FormulaIdentifier != "a" {
		t.Fatalf("expected formula identifier to be preserved, got %q", captured.Metrics[0].FormulaIdentifier)
	}
	if len(captured.Formulas) != 1 {
		t.Fatalf("expected 1 formula, got %d", len(captured.Formulas))
	}
	if captured.Formulas[0].Formula != "a" {
		t.Fatalf("expected formula to be preserved, got %q", captured.Formulas[0].Formula)
	}
	if captured.Formulas[0].Label != "" {
		t.Fatalf("expected formula label to be stripped from downstream API request, got %q", captured.Formulas[0].Label)
	}
}

func TestSanitizeFormulasLeavesExistingBehaviorUnchangedWhenLabelsAreOmitted(t *testing.T) {
	formulas := []model.Formula{
		{Formula: "a / b"},
	}

	sanitized := sanitizeFormulas(formulas)

	if len(sanitized) != 1 {
		t.Fatalf("expected 1 sanitized formula, got %d", len(sanitized))
	}
	if sanitized[0].Formula != "a / b" {
		t.Fatalf("expected formula to be preserved, got %q", sanitized[0].Formula)
	}
	if sanitized[0].Label != "" {
		t.Fatalf("expected omitted label to stay empty, got %q", sanitized[0].Label)
	}
}
