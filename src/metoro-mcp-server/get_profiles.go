package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"time"
)

func getProfilesHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	now := time.Now()
	fiveMinsAgo := now.Add(-5 * time.Minute)
	request := GetProfileRequest{
		StartTime: fiveMinsAgo.Unix(),
		EndTime:   now.Unix(),
	}

	if serviceName, ok := arguments["serviceName"].(string); ok && serviceName != "" {
		request.ServiceName = serviceName
	} else {
		return nil, fmt.Errorf("serviceName is required")
	}

	if containerNamesStr, ok := arguments["containerNames"].(string); ok && containerNamesStr != "" {
		var containerNames []string
		if err := json.Unmarshal([]byte(containerNamesStr), &containerNames); err != nil {
			return nil, fmt.Errorf("error parsing containerNames JSON: %v", err)
		}
		request.ContainerNames = containerNames
	}

	body, err := getProfilesMetoroCall(request)
	if err != nil {
		return nil, fmt.Errorf("error getting profiles: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("%s", string(body))), nil
}

func getProfilesMetoroCall(request GetProfileRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling profiles request: %v", err)
	}
	return MakeMetoroAPIRequest("POST", "profiles", bytes.NewBuffer(requestBody))
}
