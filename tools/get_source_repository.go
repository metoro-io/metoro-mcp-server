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
	ServiceName string `json:"serviceName" jsonschema:"required,description=The name of the service to get the source repository for"`

	// Optional: Environment to filter by. If not provided, all environments are considered
	Environments []string `json:"environments" jsonschema:"description=List of environments to search for the service in. If empty all environments will be considered"`

	// Required: Time configuration for the query
	TimeConfig utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get the source repository information for. You can use relative time (e.g. last 5 minutes) or absolute time range."`
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
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}

	body, err := getSourceRepositoryMetoroCall(ctx, arguments, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error getting source repository: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getSourceRepositoryMetoroCall(ctx context.Context, args GetSourceRepositoryHandlerArgs, startTime, endTime int64) ([]byte, error) {
	req := GetSourceRepositoryRequest{
		ServiceName:  args.ServiceName,
		Environments: args.Environments,
		StartTime:    startTime,
		EndTime:      endTime,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	return utils.MakeMetoroAPIRequest("POST", "source/repository", bytes.NewBuffer(reqBody), utils.GetAPIRequirementsFromRequest(ctx))
}
