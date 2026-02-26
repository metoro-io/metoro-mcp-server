package tools

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetServicesHandlerArgs struct {
	TimeConfig utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get services for. You can set relative or absolute time range using time_config"`
}

func GetServicesHandler(ctx context.Context, arguments GetServicesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	body, err := getServicesMetoroCall(ctx, arguments)
	if err != nil {
		return nil, fmt.Errorf("error getting services: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getServicesMetoroCall(ctx context.Context, arguments GetServicesHandlerArgs) ([]byte, error) {
	path, err := buildGetServicesPath(arguments)
	if err != nil {
		return nil, err
	}

	return utils.MakeMetoroAPIRequest("GET", path, nil, utils.GetAPIRequirementsFromRequest(ctx))
}

func buildGetServicesPath(arguments GetServicesHandlerArgs) (string, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return "", fmt.Errorf("error calculating time range: %v", err)
	}

	query := url.Values{}
	query.Set("startTime", strconv.FormatInt(startTime, 10))
	query.Set("endTime", strconv.FormatInt(endTime, 10))

	return "services?" + query.Encode(), nil
}
