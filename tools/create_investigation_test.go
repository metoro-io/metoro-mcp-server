package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TestCreateInvestigationHandlerRequiresCategory(t *testing.T) {
	_, err := CreateInvestigationHandler(context.Background(), CreateInvestigationHandlerArgs{
		Title:      "title",
		Category:   "",
		Summary:    "summary",
		Markdown:   "markdown",
		TimeConfig: investigationAbsoluteTimeConfig(),
	})
	if err == nil {
		t.Fatalf("expected error for missing category")
	}
	if !strings.Contains(err.Error(), "invalid category") {
		t.Fatalf("expected invalid category error, got %v", err)
	}
}

func TestCreateInvestigationHandlerAllowsDeploymentWithoutVerdict(t *testing.T) {
	var mu sync.Mutex
	var captured *model.CreateInvestigationRequest
	capturedHasVerdict := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/investigation" {
			t.Fatalf("expected path /api/v1/investigation, got %s", r.URL.Path)
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		var req model.CreateInvestigationRequest
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(bodyBytes, &raw); err != nil {
			t.Fatalf("failed to unmarshal request into raw map: %v", err)
		}

		mu.Lock()
		capturedReq := req
		captured = &capturedReq
		_, capturedHasVerdict = raw["verdict"]
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"uuid":"new-investigation"}`))
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := CreateInvestigationHandler(context.Background(), CreateInvestigationHandlerArgs{
		Title:      "title",
		Category:   investigationCategoryDeploymentVerification,
		Summary:    "summary",
		Markdown:   "markdown",
		TimeConfig: investigationAbsoluteTimeConfig(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if captured == nil {
		t.Fatalf("expected request to be captured")
	}
	if captured.Category != investigationCategoryDeploymentVerification {
		t.Fatalf("expected category %q, got %q", investigationCategoryDeploymentVerification, captured.Category)
	}
	if capturedHasVerdict {
		t.Fatalf("expected verdict field to be omitted from create payload")
	}
	if _, ok := captured.Tags["verdict"]; ok {
		t.Fatalf("expected tags.verdict to be omitted when verdict is not provided")
	}
}

func TestCreateInvestigationHandlerAcceptsDeploymentWithVerdict(t *testing.T) {
	verdict := " healthy "

	var mu sync.Mutex
	var captured *model.CreateInvestigationRequest
	capturedHasVerdict := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/investigation" {
			t.Fatalf("expected path /api/v1/investigation, got %s", r.URL.Path)
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		var req model.CreateInvestigationRequest
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(bodyBytes, &raw); err != nil {
			t.Fatalf("failed to unmarshal request into raw map: %v", err)
		}

		mu.Lock()
		capturedReq := req
		captured = &capturedReq
		_, capturedHasVerdict = raw["verdict"]
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"uuid":"new-investigation"}`))
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := CreateInvestigationHandler(context.Background(), CreateInvestigationHandlerArgs{
		Title:      "title",
		Category:   investigationCategoryDeploymentVerification,
		Verdict:    &verdict,
		Summary:    "summary",
		Markdown:   "markdown",
		TimeConfig: investigationAbsoluteTimeConfig(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if captured == nil {
		t.Fatalf("expected request to be captured")
	}
	if captured.Category != investigationCategoryDeploymentVerification {
		t.Fatalf("expected category %q, got %q", investigationCategoryDeploymentVerification, captured.Category)
	}
	if capturedHasVerdict {
		t.Fatalf("expected verdict field to be omitted from create payload")
	}
	if captured.Tags["verdict"] != "healthy" {
		t.Fatalf("expected tags.verdict to be %q, got %q", "healthy", captured.Tags["verdict"])
	}
}

func TestCreateInvestigationHandlerAddsEnvironmentAndNamespaceTags(t *testing.T) {
	environment := " production "
	namespace := " payments "

	var mu sync.Mutex
	var captured *model.CreateInvestigationRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/investigation" {
			t.Fatalf("expected path /api/v1/investigation, got %s", r.URL.Path)
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		var req model.CreateInvestigationRequest
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		mu.Lock()
		capturedReq := req
		captured = &capturedReq
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"uuid":"new-investigation"}`))
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := CreateInvestigationHandler(context.Background(), CreateInvestigationHandlerArgs{
		Title:       "title",
		Category:    investigationCategoryAnomalyInvestigation,
		Summary:     "summary",
		Markdown:    "markdown",
		Environment: &environment,
		Namespace:   &namespace,
		TimeConfig:  investigationAbsoluteTimeConfig(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if captured == nil {
		t.Fatalf("expected request to be captured")
	}
	if captured.Tags["environment"] != "production" {
		t.Fatalf("expected tags.environment to be %q, got %q", "production", captured.Tags["environment"])
	}
	if captured.Tags["namespace"] != "payments" {
		t.Fatalf("expected tags.namespace to be %q, got %q", "payments", captured.Tags["namespace"])
	}
}

func investigationAbsoluteTimeConfig() utils.TimeConfig {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"
	return utils.TimeConfig{
		Type:      utils.AbsoluteTimeRange,
		StartTime: &start,
		EndTime:   &end,
	}
}
