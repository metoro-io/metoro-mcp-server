package tools

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
)

func TestTrimLogsPayloadTextTruncatesLargeFields(t *testing.T) {
	betweenStrictAndNormal := stackTraceValueLengthLimit + 120

	logsResponse := model.GetLogsResponse{
		Logs: []model.Log{
			{
				Message: strings.Repeat("m", logMessageLengthLimit+200),
				LogAttributes: map[string]string{
					"errorVerbose":       strings.Repeat("v", stackTraceValueLengthLimit+200),
					"Error":              strings.Repeat("e", stackTraceValueLengthLimit+200),
					"stacktrace":         strings.Repeat("s", stackTraceValueLengthLimit+200),
					"user.id":            strings.Repeat("u", logAttributeValueLengthLimit+200),
					"service.error.code": strings.Repeat("c", betweenStrictAndNormal),
				},
				ResourceAttributes: map[string]string{
					"ErrorVerbose":     strings.Repeat("r", stackTraceValueLengthLimit+200),
					"my_stacktrace_v2": strings.Repeat("x", betweenStrictAndNormal),
					"k8s.pod":          strings.Repeat("p", logAttributeValueLengthLimit+200),
				},
			},
		},
	}

	raw, err := json.Marshal(logsResponse)
	if err != nil {
		t.Fatalf("failed to marshal test response: %v", err)
	}

	trimmed, changed, err := trimLogsPayloadText(string(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !changed {
		t.Fatalf("expected payload to be changed")
	}

	var parsed model.GetLogsResponse
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		t.Fatalf("failed to unmarshal trimmed payload: %v", err)
	}

	gotLog := parsed.Logs[0]
	if utf8.RuneCountInString(gotLog.Message) > logMessageLengthLimit {
		t.Fatalf("expected message to be <= %d runes, got %d", logMessageLengthLimit, utf8.RuneCountInString(gotLog.Message))
	}
	if !strings.HasSuffix(gotLog.Message, truncatedValueSuffix) {
		t.Fatalf("expected message to include truncation suffix")
	}

	strictKeys := []string{"errorVerbose", "Error", "stacktrace"}
	for _, key := range strictKeys {
		value := gotLog.LogAttributes[key]
		if utf8.RuneCountInString(value) > stackTraceValueLengthLimit {
			t.Fatalf("expected %s to be <= %d runes, got %d", key, stackTraceValueLengthLimit, utf8.RuneCountInString(value))
		}
		if !strings.HasSuffix(value, truncatedValueSuffix) {
			t.Fatalf("expected %s to include truncation suffix", key)
		}
	}

	userID := gotLog.LogAttributes["user.id"]
	if utf8.RuneCountInString(userID) > logAttributeValueLengthLimit {
		t.Fatalf("expected regular attribute to be <= %d runes, got %d", logAttributeValueLengthLimit, utf8.RuneCountInString(userID))
	}
	if !strings.HasSuffix(userID, truncatedValueSuffix) {
		t.Fatalf("expected regular oversized attribute to include truncation suffix")
	}

	controlLogAttr := gotLog.LogAttributes["service.error.code"]
	if utf8.RuneCountInString(controlLogAttr) != betweenStrictAndNormal {
		t.Fatalf("expected service.error.code to remain unchanged at %d runes, got %d", betweenStrictAndNormal, utf8.RuneCountInString(controlLogAttr))
	}

	pod := gotLog.ResourceAttributes["k8s.pod"]
	if utf8.RuneCountInString(pod) > logAttributeValueLengthLimit {
		t.Fatalf("expected resource attribute to be <= %d runes, got %d", logAttributeValueLengthLimit, utf8.RuneCountInString(pod))
	}
	if !strings.HasSuffix(pod, truncatedValueSuffix) {
		t.Fatalf("expected regular oversized resource attribute to include truncation suffix")
	}

	errorVerboseResource := gotLog.ResourceAttributes["ErrorVerbose"]
	if utf8.RuneCountInString(errorVerboseResource) > stackTraceValueLengthLimit {
		t.Fatalf("expected ErrorVerbose to be <= %d runes, got %d", stackTraceValueLengthLimit, utf8.RuneCountInString(errorVerboseResource))
	}
	if !strings.HasSuffix(errorVerboseResource, truncatedValueSuffix) {
		t.Fatalf("expected ErrorVerbose to include truncation suffix")
	}

	controlResourceAttr := gotLog.ResourceAttributes["my_stacktrace_v2"]
	if utf8.RuneCountInString(controlResourceAttr) != betweenStrictAndNormal {
		t.Fatalf("expected my_stacktrace_v2 to remain unchanged at %d runes, got %d", betweenStrictAndNormal, utf8.RuneCountInString(controlResourceAttr))
	}
}

func TestTrimLogsPayloadTextLeavesNonLogPayloadUnchanged(t *testing.T) {
	raw := `{"foo":"bar"}`
	trimmed, changed, err := trimLogsPayloadText(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if changed {
		t.Fatalf("expected no change for non-log payload")
	}
	if trimmed != raw {
		t.Fatalf("expected payload to remain unchanged")
	}
}

func TestTrimLogsPayloadTextPreservesTopLevelFields(t *testing.T) {
	rawPayload := map[string]any{
		"logs": []model.Log{
			{
				Message: strings.Repeat("m", logMessageLengthLimit+1),
			},
		},
		"cursor": "next-page",
		"metadata": map[string]any{
			"source": "context",
		},
	}

	raw, err := json.Marshal(rawPayload)
	if err != nil {
		t.Fatalf("failed to marshal test response: %v", err)
	}

	trimmed, changed, err := trimLogsPayloadText(string(raw))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !changed {
		t.Fatalf("expected payload to be changed")
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		t.Fatalf("failed to unmarshal trimmed payload: %v", err)
	}

	var cursor string
	if err := json.Unmarshal(parsed["cursor"], &cursor); err != nil {
		t.Fatalf("failed to unmarshal cursor: %v", err)
	}
	if cursor != "next-page" {
		t.Fatalf("expected cursor to be preserved, got %q", cursor)
	}

	var metadata map[string]string
	if err := json.Unmarshal(parsed["metadata"], &metadata); err != nil {
		t.Fatalf("failed to unmarshal metadata: %v", err)
	}
	if metadata["source"] != "context" {
		t.Fatalf("expected metadata source to be preserved, got %q", metadata["source"])
	}

	var logs []model.Log
	if err := json.Unmarshal(parsed["logs"], &logs); err != nil {
		t.Fatalf("failed to unmarshal logs: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected one log, got %d", len(logs))
	}
	if !strings.HasSuffix(logs[0].Message, truncatedValueSuffix) {
		t.Fatalf("expected log message to include truncation suffix")
	}
}

func TestLogsToolResponseGuardFitsOversizedLogsWithDynamicTruncation(t *testing.T) {
	const maxTokens = 900

	logsResponse := model.GetLogsResponse{
		Logs: make([]model.Log, 5),
	}
	for i := range logsResponse.Logs {
		logsResponse.Logs[i] = model.Log{
			Time:     1700000000000 + int64(i),
			Severity: "ERROR",
			Message:  strings.Repeat("message", logMessageLengthLimit),
			LogAttributes: map[string]string{
				"payload":    strings.Repeat("payload", logAttributeValueLengthLimit),
				"stacktrace": strings.Repeat("stacktrace", stackTraceValueLengthLimit),
			},
			ResourceAttributes: map[string]string{
				"k8s.pod": strings.Repeat("resource", logAttributeValueLengthLimit),
			},
			ServiceName: "checkout",
			Environment: "production",
		}
	}

	rawPayload := map[string]any{
		"logs":   logsResponse.Logs,
		"cursor": "next-page",
		"metadata": map[string]any{
			"source": "context",
		},
	}
	raw, err := json.Marshal(rawPayload)
	if err != nil {
		t.Fatalf("failed to marshal test response: %v", err)
	}

	staticOnlyResponse := mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(raw)))
	if _, err := trimLargeLogFieldsInToolResponse("get_logs", staticOnlyResponse); err != nil {
		t.Fatalf("expected static trim to succeed: %v", err)
	}
	staticTokenCount, err := estimateToolResponseTokens(staticOnlyResponse)
	if err != nil {
		t.Fatalf("expected token estimate to succeed: %v", err)
	}
	if staticTokenCount <= maxTokens {
		t.Fatalf("test fixture should still be too large after static trimming; got %d tokens", staticTokenCount)
	}

	guard := NewLogsToolResponseGuard(ToolResponseGuardOptions{MaxTokens: maxTokens})
	guarded, err := guard("get_logs", mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(raw))))
	if err != nil {
		t.Fatalf("expected guard to fit oversized logs, got %v", err)
	}
	if guarded == nil {
		t.Fatalf("expected guarded response")
	}

	tokenCount, err := estimateToolResponseTokens(guarded)
	if err != nil {
		t.Fatalf("expected token estimate to succeed: %v", err)
	}
	if tokenCount > maxTokens {
		t.Fatalf("expected guarded response to fit under %d tokens, got %d", maxTokens, tokenCount)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal([]byte(guarded.Content[0].TextContent.Text), &parsed); err != nil {
		t.Fatalf("failed to unmarshal guarded payload: %v", err)
	}

	var cursor string
	if err := json.Unmarshal(parsed["cursor"], &cursor); err != nil {
		t.Fatalf("failed to unmarshal cursor: %v", err)
	}
	if cursor != "next-page" {
		t.Fatalf("expected cursor to be preserved, got %q", cursor)
	}

	var metadata map[string]string
	if err := json.Unmarshal(parsed["metadata"], &metadata); err != nil {
		t.Fatalf("failed to unmarshal metadata: %v", err)
	}
	if metadata["source"] != "context" {
		t.Fatalf("expected metadata source to be preserved, got %q", metadata["source"])
	}

	var parsedLogs []model.Log
	if err := json.Unmarshal(parsed["logs"], &parsedLogs); err != nil {
		t.Fatalf("failed to unmarshal logs: %v", err)
	}
	if len(parsedLogs) != len(logsResponse.Logs) {
		t.Fatalf("expected %d logs to be preserved, got %d", len(logsResponse.Logs), len(parsedLogs))
	}

	dynamicallyTrimmedMessage := false
	for _, log := range parsedLogs {
		if !strings.HasSuffix(log.Message, truncatedValueSuffix) {
			t.Fatalf("expected dynamically fitted log message to include truncation suffix")
		}
		if utf8.RuneCountInString(log.Message) < logMessageLengthLimit {
			dynamicallyTrimmedMessage = true
		}

		for key, value := range log.LogAttributes {
			if !strings.HasSuffix(value, truncatedValueSuffix) {
				t.Fatalf("expected log attribute %s to include truncation suffix", key)
			}
		}
		for key, value := range log.ResourceAttributes {
			if !strings.HasSuffix(value, truncatedValueSuffix) {
				t.Fatalf("expected resource attribute %s to include truncation suffix", key)
			}
		}
	}
	if !dynamicallyTrimmedMessage {
		t.Fatalf("expected at least one message to be trimmed below the static limit")
	}
}

func TestLogsToolResponseGuardReturnsTooLargeForOversizedNonLogPayload(t *testing.T) {
	raw := `{"foo":"` + strings.Repeat("x", 4000) + `"}`
	response := mcpgolang.NewToolResponse(mcpgolang.NewTextContent(raw))

	guard := NewLogsToolResponseGuard(ToolResponseGuardOptions{MaxTokens: 50})
	guarded, err := guard("get_logs", response)
	if err == nil {
		t.Fatalf("expected oversized non-log payload to fail")
	}
	if err.Error() != toolResponseTooLargeErrorMessage {
		t.Fatalf("expected %q, got %q", toolResponseTooLargeErrorMessage, err.Error())
	}
	if guarded != nil {
		t.Fatalf("expected nil response when guard fails")
	}
	if response.Content[0].TextContent.Text != raw {
		t.Fatalf("expected non-log payload to remain unchanged")
	}
}
