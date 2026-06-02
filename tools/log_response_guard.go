package tools

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
)

const (
	logMessageLengthLimit        = 2000
	logAttributeValueLengthLimit = 600
	stackTraceValueLengthLimit   = 300
	truncatedValueSuffix         = "... [truncated]"
	minDynamicLogFieldLength     = len(truncatedValueSuffix)
)

var strictLogAttributeKeys = map[string]struct{}{
	"errorverbose": {},
	"stacktrace":   {},
	"error":        {},
}

var LogsToolResponseGuard = NewLogsToolResponseGuard(ToolResponseGuardOptions{})

func NewLogsToolResponseGuard(options ToolResponseGuardOptions) ToolResponseGuard {
	return func(toolName string, response *mcpgolang.ToolResponse) (*mcpgolang.ToolResponse, error) {
		if response == nil {
			return nil, nil
		}

		maxTokens := resolveToolResponseMaxTokens(options)
		guardedResponse, err := trimAndFitLogsToolResponse(response, maxTokens)
		if err != nil {
			return nil, err
		}

		tokenCount, err := estimateToolResponseTokens(guardedResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate response token size for tool %q: %w", toolName, err)
		}

		if tokenCount > maxTokens {
			return nil, fmt.Errorf(resolveToolResponseTooLargeMessage(options))
		}

		return guardedResponse, nil
	}
}

func trimAndFitLogsToolResponse(response *mcpgolang.ToolResponse, maxTokens int) (*mcpgolang.ToolResponse, error) {
	trimmedResponse, err := trimLargeLogFieldsInToolResponse("", response)
	if err != nil {
		return nil, err
	}

	tokenCount, err := estimateToolResponseTokens(trimmedResponse)
	if err != nil {
		return nil, err
	}
	if tokenCount <= maxTokens {
		return trimmedResponse, nil
	}

	for _, content := range trimmedResponse.Content {
		if content == nil || content.Type != mcpgolang.ContentTypeText || content.TextContent == nil {
			continue
		}

		err := fitLogsContentToToolResponseBudget(content, trimmedResponse, maxTokens)
		if err != nil {
			return nil, err
		}

		tokenCount, err = estimateToolResponseTokens(trimmedResponse)
		if err != nil {
			return nil, err
		}
		if tokenCount <= maxTokens {
			return trimmedResponse, nil
		}
	}

	return trimmedResponse, nil
}

func fitLogsContentToToolResponseBudget(content *mcpgolang.Content, response *mcpgolang.ToolResponse, maxTokens int) error {
	logsResponse, ok := parseLogsPayloadText(content.TextContent.Text)
	if !ok {
		return nil
	}

	if err := setLogsContentText(content, logsResponse); err != nil {
		return err
	}

	for {
		tokenCount, err := estimateToolResponseTokens(response)
		if err != nil {
			return err
		}
		if tokenCount <= maxTokens {
			return nil
		}

		field := findLongestShrinkableLogField(&logsResponse)
		if field == nil {
			return nil
		}

		nextLimit := utf8.RuneCountInString(field.value) / 2
		if nextLimit < minDynamicLogFieldLength {
			nextLimit = minDynamicLogFieldLength
		}

		truncated, wasTruncated := truncateWithSuffix(field.value, nextLimit)
		if !wasTruncated {
			return nil
		}

		field.set(truncated)
		if err := setLogsContentText(content, logsResponse); err != nil {
			return err
		}
	}
}

func setLogsContentText(content *mcpgolang.Content, logsResponse model.GetLogsResponse) error {
	serialized, err := json.Marshal(logsResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal trimmed logs response: %w", err)
	}
	content.TextContent.Text = string(serialized)
	return nil
}

type logStringFieldRef struct {
	value string
	set   func(string)
}

func findLongestShrinkableLogField(logsResponse *model.GetLogsResponse) *logStringFieldRef {
	var result *logStringFieldRef
	var resultLength int

	consider := func(value string, set func(string)) {
		valueLength := utf8.RuneCountInString(value)
		if valueLength <= minDynamicLogFieldLength || valueLength <= resultLength {
			return
		}

		result = &logStringFieldRef{
			value: value,
			set:   set,
		}
		resultLength = valueLength
	}

	for i := range logsResponse.Logs {
		logIndex := i
		consider(logsResponse.Logs[logIndex].Message, func(value string) {
			logsResponse.Logs[logIndex].Message = value
		})

		considerLogAttributeFields(logsResponse.Logs[logIndex].LogAttributes, consider)
		considerLogAttributeFields(logsResponse.Logs[logIndex].ResourceAttributes, consider)
	}

	return result
}

func considerLogAttributeFields(attributes map[string]string, consider func(string, func(string))) {
	keys := make([]string, 0, len(attributes))
	for key := range attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		attributeKey := key
		consider(attributes[attributeKey], func(value string) {
			attributes[attributeKey] = value
		})
	}
}

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
	logsResponse, ok := parseLogsPayloadText(raw)
	if !ok {
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

func parseLogsPayloadText(raw string) (model.GetLogsResponse, bool) {
	var rawObject map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &rawObject); err != nil {
		return model.GetLogsResponse{}, false
	}
	if _, ok := rawObject["logs"]; !ok {
		return model.GetLogsResponse{}, false
	}

	var logsResponse model.GetLogsResponse
	if err := json.Unmarshal([]byte(raw), &logsResponse); err != nil {
		return model.GetLogsResponse{}, false
	}

	return logsResponse, true
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
