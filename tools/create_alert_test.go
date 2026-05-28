package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/metoro-io/metoro-mcp-server/model"
)

func TestCreateAlertHandlerSendsIsTrainingOnPublicAlertSchema(t *testing.T) {
	tests := []struct {
		name           string
		isTraining     bool
		wantField      bool
		wantIsTraining bool
	}{
		{
			name:           "training enabled",
			isTraining:     true,
			wantField:      true,
			wantIsTraining: true,
		},
		{
			name:           "training omitted by default",
			isTraining:     false,
			wantField:      false,
			wantIsTraining: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var captured model.CreateUpdateAlertRequest
			var sawAlertUpdate bool

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/v1/metrics/attributes":
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"attributes":[]}`))
				case "/api/v1/metoroql/convert/metricSpecifierToMetoroql":
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"queries":["count(traces)"]}`))
				case "/api/v1/alerts/update":
					sawAlertUpdate = true
					if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
						t.Fatalf("failed to decode alert request: %v", err)
					}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"id":"alert-id"}`))
				default:
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
			}))
			defer server.Close()

			setMetoroAPIEnv(t, server.URL)

			_, err := CreateAlertHandler(context.Background(), CreateAlertHandlerArgs{
				AlertName:        "High latency",
				AlertDescription: "Latency is above threshold",
				Timeseries: []model.MetricSpecifier{
					{
						MetricType:        model.Trace,
						Aggregation:       model.AggregationCount,
						BucketSize:        60,
						FormulaIdentifier: "a",
					},
				},
				Formula:           model.Formula{Formula: "a"},
				Condition:         "GreaterThan",
				Threshold:         10,
				DatapointsToAlarm: 3,
				EvaluationWindow:  5,
				IsTraining:        tt.isTraining,
			})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !sawAlertUpdate {
				t.Fatalf("expected create_alert to call /api/v1/alerts/update")
			}
			if captured.Alert.HasIsTraining() != tt.wantField {
				t.Fatalf("expected isTraining field presence %v, got %v", tt.wantField, captured.Alert.HasIsTraining())
			}
			if captured.Alert.GetIsTraining() != tt.wantIsTraining {
				t.Fatalf("expected isTraining %v, got %v", tt.wantIsTraining, captured.Alert.GetIsTraining())
			}
		})
	}
}
