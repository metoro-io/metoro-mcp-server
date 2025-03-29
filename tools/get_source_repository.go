package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetSourceRepositoryHandlerArgs struct {
	// Required: Service name to get the source repository for
	ServiceName string `json:"serviceName" jsonschema:"required description=The name of the service to get the source repository for"`

	// Optional: Environment to filter by. If not provided, all environments are considered
	Environments []string `json:"environments" jsonschema:"description=List of environments to search for the service in. If empty all environments will be considered"`

	// Required: Timestamp to get metadata updates after this time
	StartTime int64 `json:"startTime" jsonschema:"required description=Unix timestamp (in milliseconds) to get the source repository information after this time"`

	// Required: Timestamp to get metadata updates before this time
	EndTime int64 `json:"endTime" jsonschema:"required description=Unix timestamp (in milliseconds) to get the source repository information before this time"`
}

type GetSourceRepositoryRequest struct {
	ServiceName  string   `json:"serviceName"`
	Environments []string `json:"environments"`
	StartTime    int64    `json:"startTime"`
	EndTime      int64    `json:"endTime"`
}

type GetSourceRepositoryResponse struct {
	// The source repository URL/path found in the deployment
	Repository string `json:"repository"`

	// Whether a repository was found
	Found bool `json:"found"`

	// The environment where the repository information was found
	Environment string `json:"environment,omitempty"`
}

func GetSourceRepositoryHandler(ctx context.Context, arguments GetSourceRepositoryHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getSourceRepositoryMetoroCall(ctx, arguments)
	if err != nil {
		return nil, fmt.Errorf("error getting source repository: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getSourceRepositoryMetoroCall(ctx context.Context, args GetSourceRepositoryHandlerArgs) ([]byte, error) {
	req := GetSourceRepositoryRequest{
		ServiceName:  args.ServiceName,
		Environments: args.Environments,
		StartTime:    args.StartTime,
		EndTime:      args.EndTime,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	return utils.MakeMetoroAPIRequest("POST", "api/v1/source/repository", bytes.NewBuffer(reqBody), utils.GetAPIRequirementsFromRequest(ctx))
}
