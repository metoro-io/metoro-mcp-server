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
)

func TestUpdateInvestigationHandlerClosingDeploymentWithoutVerdictFails(t *testing.T) {
	inProgress := false
	putCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/investigation" {
			t.Fatalf("expected path /api/v1/investigation, got %s", r.URL.Path)
		}

		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("uuid") != "inv-uuid" {
				t.Fatalf("expected uuid query param inv-uuid, got %s", r.URL.Query().Get("uuid"))
			}
			_, _ = w.Write([]byte(`{"category":"deployment_verification"}`))
		case http.MethodPut:
			putCalled = true
			_, _ = w.Write([]byte(`{"uuid":"inv-uuid"}`))
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := UpdateInvestigationHandler(context.Background(), validUpdateInvestigationArgs(func(args *UpdateInvestigationHandlerArgs) {
		args.InProgress = &inProgress
	}))
	if err == nil {
		t.Fatalf("expected error when closing deployment investigation without verdict")
	}
	if !strings.Contains(err.Error(), "verdict is required") {
		t.Fatalf("expected verdict required error, got %v", err)
	}
	if putCalled {
		t.Fatalf("expected no PUT request when validation fails")
	}
}

func TestUpdateInvestigationHandlerClosingDeploymentWithPendingFails(t *testing.T) {
	inProgress := false
	verdict := "pending"
	putCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/investigation" {
			t.Fatalf("expected path /api/v1/investigation, got %s", r.URL.Path)
		}

		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("uuid") != "inv-uuid" {
				t.Fatalf("expected uuid query param inv-uuid, got %s", r.URL.Query().Get("uuid"))
			}
			_, _ = w.Write([]byte(`{"category":"deployment_verification"}`))
		case http.MethodPut:
			putCalled = true
			_, _ = w.Write([]byte(`{"uuid":"inv-uuid"}`))
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := UpdateInvestigationHandler(context.Background(), validUpdateInvestigationArgs(func(args *UpdateInvestigationHandlerArgs) {
		args.InProgress = &inProgress
		args.Verdict = &verdict
	}))
	if err == nil {
		t.Fatalf("expected error when closing deployment investigation with pending verdict")
	}
	if !strings.Contains(err.Error(), "pending verdict is not allowed") {
		t.Fatalf("expected pending-not-allowed error, got %v", err)
	}
	if putCalled {
		t.Fatalf("expected no PUT request when validation fails")
	}
}

func TestUpdateInvestigationHandlerClosingDeploymentWithFinalVerdictSucceeds(t *testing.T) {
	inProgress := false
	validVerdicts := []string{"healthy", "degraded", "failed"}

	for _, verdictValue := range validVerdicts {
		t.Run(verdictValue, func(t *testing.T) {
			var mu sync.Mutex
			putCalled := false
			var captured *model.UpdateInvestigationRequest

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/investigation" {
					t.Fatalf("expected path /api/v1/investigation, got %s", r.URL.Path)
				}

				switch r.Method {
				case http.MethodGet:
					if r.URL.Query().Get("uuid") != "inv-uuid" {
						t.Fatalf("expected uuid query param inv-uuid, got %s", r.URL.Query().Get("uuid"))
					}
					_, _ = w.Write([]byte(`{"category":"deployment_verification"}`))
				case http.MethodPut:
					var req model.UpdateInvestigationRequest
					if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
						t.Fatalf("failed to decode request body: %v", err)
					}
					mu.Lock()
					putCalled = true
					capturedReq := req
					captured = &capturedReq
					mu.Unlock()
					_, _ = w.Write([]byte(`{"uuid":"inv-uuid"}`))
				default:
					t.Fatalf("unexpected method %s", r.Method)
				}
			}))
			defer server.Close()

			setMetoroAPIEnv(t, server.URL)

			verdict := verdictValue
			_, err := UpdateInvestigationHandler(context.Background(), validUpdateInvestigationArgs(func(args *UpdateInvestigationHandlerArgs) {
				args.InProgress = &inProgress
				args.Verdict = &verdict
			}))
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			mu.Lock()
			defer mu.Unlock()

			if !putCalled {
				t.Fatalf("expected PUT request to be sent")
			}
			if captured == nil || captured.Verdict == nil || *captured.Verdict != verdictValue {
				t.Fatalf("expected update payload verdict %q, got %v", verdictValue, captured)
			}
		})
	}
}

func TestUpdateInvestigationHandlerAddsEnvironmentAndNamespaceTags(t *testing.T) {
	environment := " production "
	namespace := " payments "

	var mu sync.Mutex
	putCalled := false
	var captured *model.UpdateInvestigationRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/investigation" {
			t.Fatalf("expected path /api/v1/investigation, got %s", r.URL.Path)
		}

		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("uuid") != "inv-uuid" {
				t.Fatalf("expected uuid query param inv-uuid, got %s", r.URL.Query().Get("uuid"))
			}
			_, _ = w.Write([]byte(`{"category":"anomaly_investigation"}`))
		case http.MethodPut:
			var req model.UpdateInvestigationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request body: %v", err)
			}
			mu.Lock()
			putCalled = true
			capturedReq := req
			captured = &capturedReq
			mu.Unlock()
			_, _ = w.Write([]byte(`{"uuid":"inv-uuid"}`))
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	setMetoroAPIEnv(t, server.URL)

	_, err := UpdateInvestigationHandler(context.Background(), validUpdateInvestigationArgs(func(args *UpdateInvestigationHandlerArgs) {
		args.Environment = &environment
		args.Namespace = &namespace
	}))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if !putCalled {
		t.Fatalf("expected PUT request to be sent")
	}
	if captured == nil || captured.Tags == nil {
		t.Fatalf("expected update payload tags to be sent, got %v", captured)
	}
	if (*captured.Tags)["environment"] != "production" {
		t.Fatalf("expected tags.environment to be %q, got %q", "production", (*captured.Tags)["environment"])
	}
	if (*captured.Tags)["namespace"] != "payments" {
		t.Fatalf("expected tags.namespace to be %q, got %q", "payments", (*captured.Tags)["namespace"])
	}
}

func validUpdateInvestigationArgs(mutate func(*UpdateInvestigationHandlerArgs)) UpdateInvestigationHandlerArgs {
	args := UpdateInvestigationHandlerArgs{
		InvestigationUUID: "inv-uuid",
		Title:             "title",
		Summary:           "summary",
		Markdown:          "markdown",
		TimeConfig:        investigationAbsoluteTimeConfig(),
	}
	if mutate != nil {
		mutate(&args)
	}
	return args
}
