package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-mcp-server/model"
	"github.com/metoro-mcp-server/utils"
	"time"
)

type GetMetricAttributesHandlerArgs struct {
	MetricName       string              `json:"metricName" jsonschema:"required,description=The name of the metric to get attributes for"`
	FilterAttributes map[string][]string `json:"filterAttributes" jsonschema:"description=The attributes to filter the metric attributes by"`
}

func GetMetricAttributesHandler(arguments GetMetricAttributesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := model.MetricAttributesRequest{
		StartTime:        fiveMinsAgo.Unix(),
		EndTime:          now.Unix(),
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
