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

type GetPodsHandlerArgs struct {
	ServiceName  string   `json:"serviceName" jsonschema:"description=The name of the service to get pods for. One of serviceName or nodeName is required"`
	NodeName     string   `json:"nodeName" jsonschema:"description=The name of the node to get pods for. One of serviceName or nodeName is required"`
	Environments []string `json:"environments" jsonschema:"description=The environments to get pods for. If empty, all environments will be used."`
}

func GetPodsHandler(arguments GetPodsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)

	// One of serviceName or nodeName is required.
	if arguments.ServiceName == "" && arguments.NodeName == "" {
		return nil, fmt.Errorf("one of serviceName or nodeName is required")
	}

	request := model.GetPodsRequest{
		StartTime:    fiveMinsAgo.Unix(),
		EndTime:      now.Unix(),
		Environments: arguments.Environments,
		ServiceName:  arguments.ServiceName,
		NodeName:     arguments.NodeName,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/pods", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
