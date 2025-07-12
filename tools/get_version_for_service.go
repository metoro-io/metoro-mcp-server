package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
	"gopkg.in/yaml.v3"
)

type GetVersionForServiceHandlerArgs struct {
	TimeConfig   utils.TimeConfig `json:"time_config" jsonschema:"required,description=The time to get container versions. e.g. if you want to see the versions 5 minutes ago you would set time_period=5 and time_window=Minutes. You can also set an absolute time range by setting start_time and end_time"`
	ServiceName  string           `json:"serviceName" jsonschema:"required,description=The name of the service to get container versions for."`
	Environments []string         `json:"environments" jsonschema:"description=The environments to get service versions for. If empty all environments will be used."`
}

type GetVersionForServiceResponse struct {
	ContainerVersions map[string]map[string]string `json:"container_versions"`
}

type K8sResourceSummaryResponse struct {
	K8sResourceSummary []struct {
		Environment  string `json:"environment"`
		Kind         string `json:"kind"`
		ResourceYaml string `json:"resourceYaml"`
	} `json:"k8sResourceSummary"`
}

func GetVersionForServiceHandler(ctx context.Context, arguments GetVersionForServiceHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetPodsRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		ServiceName:  arguments.ServiceName,
		Environments: arguments.Environments,
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/summary", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	// Parse the JSON response
	var summaryResponse K8sResourceSummaryResponse
	err = json.Unmarshal(resp, &summaryResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Extract container versions from each environment
	containerVersions := make(map[string]map[string]string)

	for _, resource := range summaryResponse.K8sResourceSummary {
		// Parse the YAML for each resource
		var yamlData map[string]interface{}
		err = yaml.Unmarshal([]byte(resource.ResourceYaml), &yamlData)
		if err != nil {
			continue // Skip if we can't parse this resource
		}

		// Extract containers for this environment
		envContainers := make(map[string]string)

		// Check for spec.template.spec.containers (Deployment/StatefulSet)
		if spec, ok := yamlData["spec"].(map[string]interface{}); ok {
			if template, ok := spec["template"].(map[string]interface{}); ok {
				if templateSpec, ok := template["spec"].(map[string]interface{}); ok {
					extractContainers(templateSpec, envContainers)
				}
			}
			// Also check spec.containers directly (DaemonSet)
			extractContainers(spec, envContainers)
		}

		if len(envContainers) > 0 {
			containerVersions[resource.Environment] = envContainers
		}
	}

	response := GetVersionForServiceResponse{
		ContainerVersions: containerVersions,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(string(jsonResponse))), nil
}

func extractContainers(spec map[string]interface{}, containerVersions map[string]string) {
	if containers, ok := spec["containers"].([]interface{}); ok {
		for _, container := range containers {
			if containerMap, ok := container.(map[string]interface{}); ok {
				containerName, nameOk := containerMap["name"].(string)
				image, imageOk := containerMap["image"].(string)
				if nameOk && imageOk {
					containerVersions[containerName] = image
				}
			}
		}
	}

	// Also check for init containers
	if initContainers, ok := spec["initContainers"].([]interface{}); ok {
		for _, container := range initContainers {
			if containerMap, ok := container.(map[string]interface{}); ok {
				containerName, nameOk := containerMap["name"].(string)
				image, imageOk := containerMap["image"].(string)
				if nameOk && imageOk {
					containerVersions["init-"+containerName] = image
				}
			}
		}
	}
}
