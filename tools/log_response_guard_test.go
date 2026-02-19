package tools

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

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
