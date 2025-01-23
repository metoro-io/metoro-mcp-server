package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetKubernetesObjectsHandlerArgs struct {
	// The type of kubernetes object. E.g. Pod, Node, CronJob etc
	Kind string `json:"kind" jsonschema:"required,description=The type of kubernetes object to get (e.g. Pod, Node, CronJob)"`
	// The optional namespace of the kubernetes object
	Namespace *string `json:"namespace" jsonschema:"description=The namespace of the kubernetes objects to get. If not provided, objects from all namespaces will be returned"`
	// The optional name of the kubernetes object
	Name *string `json:"name" jsonschema:"description=The name of the kubernetes object to get. If not provided, all objects of the specified kind will be returned"`
	// The optional labels of the kubernetes object
	Labels *map[string]string `json:"labels" jsonschema:"description=Labels to filter the kubernetes objects by"`
	// The optional service of the kubernetes object
	Service *string `json:"service" jsonschema:"description=The service name to filter the kubernetes objects by"`
	// The optional environment of the kubernetes object
	Environment *string `json:"environment" jsonschema:"description=The environment to get the kubernetes objects from"`
	// Time - Unix timestamp at which we should get the kubernetes objects.
	// We will return object as it existed at this point, not how it changed since
	Time *int64 `json:"time" jsonschema:"description=Unix timestamp at which to get the kubernetes objects. If not provided, current time will be used"`
}

type GetKubernetesObjectsRequest struct {
	// Time - Unix timestamp at which we should get the kubernetes objects.
	// We will return object as it existed at this point, not how it changed since
	Time int64
	// The type of kubernetes object. E.g. Pod, Node, CronJob etc
	Kind string
	// The optional namespace of the kubernetes object
	Namespace *string
	// The optional name of the kubernetes object
	Name *string
	// The optional labels of the kubernetes object
	Labels *map[string]string
	// The optional service of the kubernetes object
	Service *string
	// The optional environment of the kubernetes object
	Environment *string
}

type MetoroKubernetesObject struct {
	Kind        string          `json:"kind"`
	Environment string          `json:"environment"`
	ServiceName *string         `json:"serviceName"`
	Object      json.RawMessage `json:"object"`
}

type GetKubernetesObjectsResponse struct {
	// Returns the raw json of the kubernetes objects, can be deserialized into their underlying types
	Objects []MetoroKubernetesObject `json:"objects"`
}

func GetKubernetesObjectsHandler(ctx context.Context, arguments GetKubernetesObjectsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	requestTime := time.Now().Unix()
	if arguments.Time != nil {
		requestTime = *arguments.Time
	}

	request := GetKubernetesObjectsRequest{
		Time:        requestTime,
		Kind:        arguments.Kind,
		Namespace:   arguments.Namespace,
		Name:        arguments.Name,
		Labels:      arguments.Labels,
		Service:     arguments.Service,
		Environment: arguments.Environment,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "kubernetesObjects", bytes.NewBuffer(jsonBody), utils.GetAPIRequirementsFromRequest(ctx))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
