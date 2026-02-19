package tools

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/metoro-io/metoro-mcp-server/model"
)

func TestTrimLogsPayloadTextTruncatesLargeFields(t *testing.T) {
	logsResponse := model.GetLogsResponse{
		Logs: []model.Log{
			{
				Message: strings.Repeat("m", logMessageLengthLimit+200),
				LogAttributes: map[string]string{
					"exception.stacktrace": strings.Repeat("s", stackTraceValueLengthLimit+200),
					"user.id":              strings.Repeat("u", logAttributeValueLengthLimit+200),
				},
				ResourceAttributes: map[string]string{
					"k8s.pod": strings.Repeat("p", logAttributeValueLengthLimit+200),
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

	stacktrace := gotLog.LogAttributes["exception.stacktrace"]
	if utf8.RuneCountInString(stacktrace) > stackTraceValueLengthLimit {
		t.Fatalf("expected stacktrace to be <= %d runes, got %d", stackTraceValueLengthLimit, utf8.RuneCountInString(stacktrace))
	}
	if !strings.HasSuffix(stacktrace, truncatedValueSuffix) {
		t.Fatalf("expected stacktrace to include truncation suffix")
	}

	userID := gotLog.LogAttributes["user.id"]
	if utf8.RuneCountInString(userID) > logAttributeValueLengthLimit {
		t.Fatalf("expected regular attribute to be <= %d runes, got %d", logAttributeValueLengthLimit, utf8.RuneCountInString(userID))
	}

	pod := gotLog.ResourceAttributes["k8s.pod"]
	if utf8.RuneCountInString(pod) > logAttributeValueLengthLimit {
		t.Fatalf("expected resource attribute to be <= %d runes, got %d", logAttributeValueLengthLimit, utf8.RuneCountInString(pod))
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
