package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetMetricAttributesHandlerArgs struct {
	TimeConfig       utils.TimeConfig    `json:"timeConfig" jsonschema:"required,description=The time period to get the possible values of metric attributes"`
	MetricName       string              `json:"metricName" jsonschema:"required,description=The name of the metric to get attributes for"`
	FilterAttributes map[string][]string `json:"filterAttributes" jsonschema:"description=The attributes to filter the metric attributes by"`
}

func GetMetricAttributesHandler(arguments GetMetricAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime := utils.CalculateTimeRange(arguments.TimeConfig)
	request := model.MetricAttributesRequest{
		StartTime:        startTime,
		EndTime:          endTime,
		MetricName:       arguments.MetricName,
		FilterAttributes: arguments.FilterAttributes,
	}
	response, err := getMetricAttributesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}

func getMetricAttributesMetoroCall(request model.MetricAttributesRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return utils.MakeMetoroAPIRequest("POST", "metricAttributes", bytes.NewBuffer(jsonData))
}
