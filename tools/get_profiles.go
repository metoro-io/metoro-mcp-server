package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/model"
	"github/metoro-io/metoro-mcp-server/src/metoro-mcp-server/utils"
	"time"
)

type GetProfileHandlerArgs struct {
	ServiceName    string   `json:"serviceName" jsonschema:"required,description=The name of the service to get profiles for"`
	ContainerNames []string `json:"containerNames" jsonschema:"description=The container names to get profiles for"`
}

func GetProfilesHandler(arguments GetProfileHandlerArgs) (*mcpgolang.ToolResponse, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := model.GetProfileRequest{
		StartTime:      fiveMinsAgo.Unix(),
		EndTime:        now.Unix(),
		ServiceName:    arguments.ServiceName,
		ContainerNames: arguments.ContainerNames,
	}

	body, err := getProfilesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting profiles: %v", err)
	}
	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func getProfilesMetoroCall(request model.GetProfileRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling profiles request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "profiles", bytes.NewBuffer(requestBody))
}
