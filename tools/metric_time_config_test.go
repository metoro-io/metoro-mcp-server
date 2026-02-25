package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TestGetMetricNamesHandlerUsesAbsoluteTimeConfig(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"
	fuzzyMatch := "cpu"
	expectedStart := mustParseRFC3339Unix(t, start)
	expectedEnd := mustParseRFC3339Unix(t, end)

	var mu sync.Mutex
	var captured *model.FuzzyMetricsRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/fuzzyMetricsNames" {
			t.Fatalf("expected path /api/v1/fuzzyMetricsNames, got %s", r.URL.Path)
		}

		var req model.FuzzyMetricsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		mu.Lock()
		copied := req
		captured = &copied
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"metrics":["container.cpu.usage"]}`))
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := GetMetricNamesHandler(context.Background(), GetMetricNamesHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		FuzzyStringMatch: fuzzyMatch,
		Environments:     []string{"prod", "staging"},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if captured == nil {
		t.Fatalf("expected fuzzyMetricsNames request to be captured")
	}
	if captured.StartTime != expectedStart {
		t.Fatalf("expected startTime %d, got %d", expectedStart, captured.StartTime)
	}
	if captured.EndTime != expectedEnd {
		t.Fatalf("expected endTime %d, got %d", expectedEnd, captured.EndTime)
	}
	if captured.MetricFuzzyMatch != fuzzyMatch {
		t.Fatalf("expected metric fuzzy match %q, got %q", fuzzyMatch, captured.MetricFuzzyMatch)
	}
	expectedEnvironments := []string{"prod", "staging"}
	if strings.Join(captured.Environments, ",") != strings.Join(expectedEnvironments, ",") {
		t.Fatalf("expected environments %v, got %v", expectedEnvironments, captured.Environments)
	}
}

func TestGetMetricNamesHandlerOmitsFuzzyMatchByDefault(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	var mu sync.Mutex
	var captured *model.FuzzyMetricsRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/fuzzyMetricsNames" {
			t.Fatalf("expected path /api/v1/fuzzyMetricsNames, got %s", r.URL.Path)
		}

		var req model.FuzzyMetricsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		mu.Lock()
		copied := req
		captured = &copied
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"metrics":["container.cpu.usage"]}`))
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := GetMetricNamesHandler(context.Background(), GetMetricNamesHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if captured == nil {
		t.Fatalf("expected fuzzyMetricsNames request to be captured")
	}
	if captured.MetricFuzzyMatch != "" {
		t.Fatalf("expected empty metric fuzzy match, got %q", captured.MetricFuzzyMatch)
	}
}

func TestGetMetricNamesHandlerReturnsErrorForInvalidTimeConfig(t *testing.T) {
	_, err := GetMetricNamesHandler(context.Background(), GetMetricNamesHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type: utils.RelativeTimeRange,
		},
	})
	if err == nil {
		t.Fatalf("expected error for invalid time config")
	}
	if !strings.Contains(err.Error(), "error calculating time range") {
		t.Fatalf("expected error to contain calculation failure, got %v", err)
	}
}

func TestGetAttributeKeysHandlerUsesProvidedTimeRangeForMetricValidation(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"
	expectedStart := mustParseRFC3339Unix(t, start)
	expectedEnd := mustParseRFC3339Unix(t, end)

	var mu sync.Mutex
	var fuzzyRequest *model.FuzzyMetricsRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/fuzzyMetricsNames":
			var req model.FuzzyMetricsRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode fuzzy request: %v", err)
			}

			mu.Lock()
			copied := req
			fuzzyRequest = &copied
			mu.Unlock()

			_, _ = w.Write([]byte(`{"metrics":["container.cpu.usage"]}`))
		case "/api/v1/metrics/attributes":
			_, _ = w.Write([]byte(`{"attributes":[]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := GetAttributeKeysHandler(context.Background(), GetAttributeKeysHandlerArgs{
		Type: model.Metric,
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		MetricName: "container.cpu.usage",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if fuzzyRequest == nil {
		t.Fatalf("expected CheckMetric to call fuzzyMetricsNames")
	}
	if fuzzyRequest.StartTime != expectedStart {
		t.Fatalf("expected startTime %d, got %d", expectedStart, fuzzyRequest.StartTime)
	}
	if fuzzyRequest.EndTime != expectedEnd {
		t.Fatalf("expected endTime %d, got %d", expectedEnd, fuzzyRequest.EndTime)
	}
}

func TestCheckTimeseriesUsesProvidedTimeRangeForMetricValidation(t *testing.T) {
	expectedStart := int64(1739968800)
	expectedEnd := int64(1739969100)

	var mu sync.Mutex
	var fuzzyRequest *model.FuzzyMetricsRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/fuzzyMetricsNames":
			var req model.FuzzyMetricsRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode fuzzy request: %v", err)
			}

			mu.Lock()
			copied := req
			fuzzyRequest = &copied
			mu.Unlock()

			_, _ = w.Write([]byte(`{"metrics":["container.cpu.usage"]}`))
		case "/api/v1/metrics/attributes":
			_, _ = w.Write([]byte(`{"attributes":[]}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	err := checkTimeseries(context.Background(), []model.SingleTimeseriesRequest{
		{
			Type:       model.Metric,
			MetricName: "container.cpu.usage",
		},
	}, expectedStart, expectedEnd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if fuzzyRequest == nil {
		t.Fatalf("expected CheckMetric to call fuzzyMetricsNames")
	}
	if fuzzyRequest.StartTime != expectedStart {
		t.Fatalf("expected startTime %d, got %d", expectedStart, fuzzyRequest.StartTime)
	}
	if fuzzyRequest.EndTime != expectedEnd {
		t.Fatalf("expected endTime %d, got %d", expectedEnd, fuzzyRequest.EndTime)
	}
}

func setMetoroAPIEnv(t *testing.T, apiURL string) {
	t.Helper()

	oldURL, hadURL := os.LookupEnv(utils.METORO_API_URL_ENV_VAR)
	oldToken, hadToken := os.LookupEnv(utils.METORO_AUTH_TOKEN_ENV_VAR)

	if err := os.Setenv(utils.METORO_API_URL_ENV_VAR, apiURL); err != nil {
		t.Fatalf("failed to set %s: %v", utils.METORO_API_URL_ENV_VAR, err)
	}
	if err := os.Setenv(utils.METORO_AUTH_TOKEN_ENV_VAR, "test-token"); err != nil {
		t.Fatalf("failed to set %s: %v", utils.METORO_AUTH_TOKEN_ENV_VAR, err)
	}

	t.Cleanup(func() {
		if hadURL {
			_ = os.Setenv(utils.METORO_API_URL_ENV_VAR, oldURL)
		} else {
			_ = os.Unsetenv(utils.METORO_API_URL_ENV_VAR)
		}

		if hadToken {
			_ = os.Setenv(utils.METORO_AUTH_TOKEN_ENV_VAR, oldToken)
		} else {
			_ = os.Unsetenv(utils.METORO_AUTH_TOKEN_ENV_VAR)
		}
	})
}

func mustParseRFC3339Unix(t *testing.T, value string) int64 {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", value, err)
	}
	return parsed.Unix()
}
