package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetLogContextHandlerArgs struct {
	ExistingLogTime int64   `json:"existingLogTime" jsonschema:"required,description=The time of the existing log line to get context for in milliseconds since epoch"`
	ContainerId     string  `json:"containerId" jsonschema:"required,description=The ID of the container to get logs for"`
	ServiceName     string  `json:"serviceName" jsonschema:"required,description=The name of the service to get logs for"`
	NumLinesBefore  int     `json:"numLinesBefore" jsonschema:"required,description=Number of log lines before the existing log line to return"`
	NumLinesAfter   int     `json:"numLinesAfter" jsonschema:"required,description=Number of log lines after the existing log line to return"`
	Environment     *string `json:"environment,omitempty" jsonschema:"description=The environment to get logs for. If empty all environments will be included. Specifying this improves performance"`
}

type GetContainerContextLogsRequest struct {
	ExistingLogTime int64   `json:"existingLogTime"`
	ContainerId     string  `json:"containerId"`
	ServiceName     string  `json:"serviceName"`
	NumLinesBefore  int     `json:"numLinesBefore"`
	NumLinesAfter   int     `json:"numLinesAfter"`
	Environment     *string `json:"environment,omitempty"`
}

func GetLogContextHandler(ctx context.Context, arguments GetLogContextHandlerArgs) (*mcpgolang.ToolResponse, error) {
	request := GetContainerContextLogsRequest{
		ExistingLogTime: arguments.ExistingLogTime,
		ContainerId:     arguments.ContainerId,
		ServiceName:     arguments.ServiceName,
		NumLinesBefore:  arguments.NumLinesBefore,
		NumLinesAfter:   arguments.NumLinesAfter,
		Environment:     arguments.Environment,
	}

	resp, err := getLogContextMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting log context: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(resp))), nil
}

func getLogContextMetoroCall(ctx context.Context, request GetContainerContextLogsRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling log context request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "logs/container/context", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
