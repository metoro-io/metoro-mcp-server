package tools

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetNamespacesHandlerArgs struct {
	TimeConfig utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get namespaces for. You can set relative or absolute time range using time_config"`
}

func GetNamespacesHandler(ctx context.Context, arguments GetNamespacesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getNamespacesMetoroCall(ctx, arguments)
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getNamespacesMetoroCall(ctx context.Context, arguments GetNamespacesHandlerArgs) ([]byte, error) {
	path, err := buildGetNamespacesPath(arguments)
	if err != nil {
		return nil, err
	}

	return utils.MakeMetoroAPIRequest("GET", path, nil, utils.GetAPIRequirementsFromRequest(ctx))
}

func buildGetNamespacesPath(arguments GetNamespacesHandlerArgs) (string, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return "", fmt.Errorf("error calculating time range: %v", err)
	}

	query := url.Values{}
	query.Set("startTime", strconv.FormatInt(startTime, 10))
	query.Set("endTime", strconv.FormatInt(endTime, 10))

	return "namespaces?" + query.Encode(), nil
}
