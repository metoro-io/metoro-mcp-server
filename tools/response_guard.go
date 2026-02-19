package tools

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	mcpgolang "github.com/metoro-io/mcp-golang"
)

const (
	toolResponseMaxTokensEnvVar      = "METORO_TOOL_RESPONSE_MAX_TOKENS"
	defaultToolResponseMaxTokens     = 12000
	toolResponseTooLargeErrorMessage = "Tool response was too large to return, try more filtering/smaller time window or a different tool"
)

type ToolResponseGuard func(toolName string, response *mcpgolang.ToolResponse) (*mcpgolang.ToolResponse, error)

type ToolResponseModifier func(toolName string, response *mcpgolang.ToolResponse) (*mcpgolang.ToolResponse, error)

type ToolResponseGuardOptions struct {
	MaxTokens            int
	TooLargeErrorMessage string
}

var DefaultToolResponseGuard = NewToolResponseGuard(nil, ToolResponseGuardOptions{})

func (tool MetoroTools) WrappedHandler() any {
	return wrapToolHandlerWithResponseGuard(tool.Name, tool.Handler, tool.ResponseGuard)
}

func NewToolResponseGuard(modifier ToolResponseModifier, options ToolResponseGuardOptions) ToolResponseGuard {
	return func(toolName string, response *mcpgolang.ToolResponse) (*mcpgolang.ToolResponse, error) {
		if response == nil {
			return nil, nil
		}

		guardedResponse := response
		if modifier != nil {
			var err error
			guardedResponse, err = modifier(toolName, guardedResponse)
			if err != nil {
				return nil, err
			}
			if guardedResponse == nil {
				return nil, fmt.Errorf("tool response modifier returned nil response")
			}
		}

		maxTokens := options.MaxTokens
		if maxTokens <= 0 {
			maxTokens = getGlobalToolResponseMaxTokens()
		}

		tooLargeMessage := options.TooLargeErrorMessage
		if tooLargeMessage == "" {
			tooLargeMessage = toolResponseTooLargeErrorMessage
		}

		tokenCount, err := estimateToolResponseTokens(guardedResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate response token size for tool %q: %w", toolName, err)
		}

		if tokenCount > maxTokens {
			return nil, fmt.Errorf(tooLargeMessage)
		}

		return guardedResponse, nil
	}
}

func getGlobalToolResponseMaxTokens() int {
	value := strings.TrimSpace(os.Getenv(toolResponseMaxTokensEnvVar))
	if value == "" {
		return defaultToolResponseMaxTokens
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return defaultToolResponseMaxTokens
	}

	return parsed
}

func estimateToolResponseTokens(response *mcpgolang.ToolResponse) (int, error) {
	body, err := json.Marshal(response)
	if err != nil {
		return 0, err
	}

	return estimateTokenCount(string(body)), nil
}

func estimateTokenCount(value string) int {
	if value == "" {
		return 0
	}

	runeCount := utf8.RuneCountInString(value)
	return int(math.Ceil(float64(runeCount) / 4.0))
}

func wrapToolHandlerWithResponseGuard(toolName string, handler any, overrideGuard ToolResponseGuard) any {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	wrapped := reflect.MakeFunc(handlerType, func(inputs []reflect.Value) []reflect.Value {
		outputs := handlerValue.Call(inputs)
		if len(outputs) != 2 {
			return outputs
		}

		if !outputs[1].IsNil() || outputs[0].IsNil() {
			return outputs
		}

		response, ok := outputs[0].Interface().(*mcpgolang.ToolResponse)
		if !ok {
			outputs[0] = reflect.Zero(handlerType.Out(0))
			outputs[1] = reflect.ValueOf(fmt.Errorf("tool handler returned unexpected response type"))
			return outputs
		}

		guard := DefaultToolResponseGuard
		if overrideGuard != nil {
			guard = overrideGuard
		}

		guardedResponse, err := guard(toolName, response)
		if err != nil {
			outputs[0] = reflect.Zero(handlerType.Out(0))
			outputs[1] = reflect.ValueOf(err)
			return outputs
		}
		if guardedResponse == nil {
			outputs[0] = reflect.Zero(handlerType.Out(0))
			outputs[1] = reflect.ValueOf(fmt.Errorf("tool response guard returned nil response"))
			return outputs
		}

		outputs[0] = reflect.ValueOf(guardedResponse)
		return outputs
	})

	return wrapped.Interface()
}
