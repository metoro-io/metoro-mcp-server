package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
)

const (
	logMessageLengthLimit        = 2000
	logAttributeValueLengthLimit = 600
	stackTraceValueLengthLimit   = 300
	truncatedValueSuffix         = "... [truncated]"
)

var strictLogAttributeKeys = map[string]struct{}{
	"errorverbose": {},
	"stacktrace":   {},
	"error":        {},
}

var LogsToolResponseGuard = NewToolResponseGuard(trimLargeLogFieldsInToolResponse, ToolResponseGuardOptions{})

func trimLargeLogFieldsInToolResponse(_ string, response *mcpgolang.ToolResponse) (*mcpgolang.ToolResponse, error) {
	for _, content := range response.Content {
		if content == nil || content.Type != mcpgolang.ContentTypeText || content.TextContent == nil {
			continue
		}

		trimmedText, changed, err := trimLogsPayloadText(content.TextContent.Text)
		if err != nil {
			return nil, err
		}
		if changed {
			content.TextContent.Text = trimmedText
		}
	}

	return response, nil
}

func trimLogsPayloadText(raw string) (string, bool, error) {
	var logsResponse model.GetLogsResponse
	if err := json.Unmarshal([]byte(raw), &logsResponse); err != nil {
		// Not all tools with this guard necessarily return log payloads all the time.
		return raw, false, nil
	}

	changed := trimLogsModelResponse(&logsResponse)
	if !changed {
		return raw, false, nil
	}

	serialized, err := json.Marshal(logsResponse)
	if err != nil {
		return "", false, fmt.Errorf("failed to marshal trimmed logs response: %w", err)
	}

	return string(serialized), true, nil
}

func trimLogsModelResponse(logsResponse *model.GetLogsResponse) bool {
	changed := false
	for i := range logsResponse.Logs {
		truncatedMessage, wasTruncated := truncateWithSuffix(logsResponse.Logs[i].Message, logMessageLengthLimit)
		if wasTruncated {
			logsResponse.Logs[i].Message = truncatedMessage
			changed = true
		}

		if trimLogAttributeValues(logsResponse.Logs[i].LogAttributes) {
			changed = true
		}
		if trimLogAttributeValues(logsResponse.Logs[i].ResourceAttributes) {
			changed = true
		}
	}

	return changed
}

func trimLogAttributeValues(attributes map[string]string) bool {
	changed := false
	for key, value := range attributes {
		limit := logAttributeValueLengthLimit
		if isStrictlyTrimmedLogAttribute(key) {
			limit = stackTraceValueLengthLimit
		}

		trimmed, wasTruncated := truncateWithSuffix(value, limit)
		if wasTruncated {
			attributes[key] = trimmed
			changed = true
		}
	}

	return changed
}

func isStrictlyTrimmedLogAttribute(attributeKey string) bool {
	key := strings.ToLower(strings.TrimSpace(attributeKey))
	_, ok := strictLogAttributeKeys[key]
	return ok
}

func truncateWithSuffix(value string, limit int) (string, bool) {
	runes := []rune(value)
	if limit <= 0 || len(runes) <= limit {
		return value, false
	}

	suffixRunes := []rune(truncatedValueSuffix)
	if limit <= len(suffixRunes) {
		return string(suffixRunes[:limit]), true
	}

	return string(runes[:limit-len(suffixRunes)]) + truncatedValueSuffix, true
}
