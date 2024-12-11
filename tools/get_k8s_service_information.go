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

type GetK8sSErviceInformationHandlerArgs struct {
	ServiceName  string   `json:"serviceName" jsonschema:"required,description=The name of the service to get information for"`
	Environments []string `json:"environments" jsonschema:"description=The environments to get information for. If empty, all environments will be used."`
}

func GetK8sServiceInformationHandler(arguments GetK8sSErviceInformationHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := model.GetPodsRequest{
		StartTime:    fiveMinsAgo.Unix(),
		EndTime:      now.Unix(),
		ServiceName:  arguments.ServiceName,
		Environments: arguments.Environments,
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/summary", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
