package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
	"time"
)

type GetMetricNamesHandlerArgs struct {
	Environments []string `json:"environments" jsonschema:"description=Environments to get metrics from. If empty, all environments will be used."`
}

func GetMetricNamesHandler(arguments GetMetricNamesHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := model.FuzzyMetricsRequest{
		StartTime:        fiveMinsAgo.Unix(),
		EndTime:          now.Unix(),
		MetricFuzzyMatch: "", // This will return all the metric names.
		Environments:     arguments.Environments,
	}
	response, err := getMetricNamesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error calling Metoro API: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(response)))), nil
}

func getMetricNamesMetoroCall(request model.FuzzyMetricsRequest) ([]byte, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return utils.MakeMetoroAPIRequest("POST", "fuzzyMetricsNames", bytes.NewBuffer(jsonData))
}
