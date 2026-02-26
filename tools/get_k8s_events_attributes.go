package tools

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sEventsAttributesHandlerArgs struct {
	TimeConfig utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time period to get Kubernetes event attributes for. You can set relative or absolute time range using time_config"`
}

func GetK8sEventsAttributesHandler(ctx context.Context, arguments GetK8sEventsAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	path, err := buildGetK8sEventsSummaryAttributesPath(arguments)
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeMetoroAPIRequest("GET", path, nil, utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, err
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func buildGetK8sEventsSummaryAttributesPath(arguments GetK8sEventsAttributesHandlerArgs) (string, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return "", fmt.Errorf("error calculating time range: %v", err)
	}

	query := url.Values{}
	query.Set("startTime", strconv.FormatInt(startTime, 10))
	query.Set("endTime", strconv.FormatInt(endTime, 10))

	return "k8s/events/summaryAttributes?" + query.Encode(), nil
}
