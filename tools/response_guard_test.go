package tools

import (
	"context"
	"strings"
	"testing"

	mcpgolang "github.com/metoro-io/mcp-golang"
)

func TestWrappedHandlerReturnsTooLargeError(t *testing.T) {
	type testArgs struct {
		Query string `json:"query"`
	}

	handler := func(ctx context.Context, arguments testArgs) (*mcpgolang.ToolResponse, error) {
		return mcpgolang.NewToolResponse(
			mcpgolang.NewTextContent(strings.Repeat("x", 4000)),
		), nil
	}

	tool := MetoroTools{
		Name:    "test_tool",
		Handler: handler,
		ResponseGuard: NewToolResponseGuard(nil, ToolResponseGuardOptions{
			MaxTokens: 50,
		}),
	}

	wrappedAny := tool.WrappedHandler()
	wrapped, ok := wrappedAny.(func(context.Context, testArgs) (*mcpgolang.ToolResponse, error))
	if !ok {
		t.Fatalf("wrapped handler has unexpected type %T", wrappedAny)
	}

	response, err := wrapped(context.Background(), testArgs{})
	if err == nil {
		t.Fatalf("expected tool size error but got nil")
	}
	if err.Error() != toolResponseTooLargeErrorMessage {
		t.Fatalf("expected %q, got %q", toolResponseTooLargeErrorMessage, err.Error())
	}
	if response != nil {
		t.Fatalf("expected nil response when guard fails")
	}
}

func TestNewToolResponseGuardAppliesModifierBeforeTokenCheck(t *testing.T) {
	guard := NewToolResponseGuard(func(_ string, response *mcpgolang.ToolResponse) (*mcpgolang.ToolResponse, error) {
		response.Content[0].TextContent.Text = "small"
		return response, nil
	}, ToolResponseGuardOptions{
		MaxTokens: 20,
	})

	response := mcpgolang.NewToolResponse(
		mcpgolang.NewTextContent(strings.Repeat("x", 4000)),
	)

	guarded, err := guard("test_tool", response)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if guarded == nil {
		t.Fatalf("expected non-nil guarded response")
	}
	if guarded.Content[0].TextContent.Text != "small" {
		t.Fatalf("expected modifier to rewrite response text")
	}
}
