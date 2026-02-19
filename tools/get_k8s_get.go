package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sGetHandlerArgs struct {
	TimeConfig         utils.TimeConfig `json:"time_config" jsonschema:"required,description=Time settings for this query. End time is used as the point in time"`
	Environment        string           `json:"environment" jsonschema:"description=Optional environment filter for this query"`
	Namespace          string           `json:"namespace" jsonschema:"description=Optional namespace filter for this query"`
	ResourceAPIVersion string           `json:"resource_api_version" jsonschema:"required,description=API version of the kubernetes resource such as v1 or apps v1"`
	ResourceKind       string           `json:"resource_kind" jsonschema:"required,description=Kind of the kubernetes resource such as Pod Deployment StatefulSet"`
	Name               string           `json:"name" jsonschema:"required,description=Resource name to fetch"`
	UID                string           `json:"uid" jsonschema:"description=Optional resource uid to disambiguate when names are reused"`
	Format             string           `json:"format" jsonschema:"enum=yaml,enum=json,description=Optional response format. Use yaml or json"`
}

type GetK8sGetRequest struct {
	Environment *string              `json:"environment,omitempty"`
	Namespace   *string              `json:"namespace,omitempty"`
	Resource    k8sResourceReference `json:"resource"`
	Name        string               `json:"name"`
	UID         *string              `json:"uid,omitempty"`
	Time        *int64               `json:"time"`
	Format      *string              `json:"format,omitempty"`
}

func GetK8sGetHandler(ctx context.Context, arguments GetK8sGetHandlerArgs) (*mcpgolang.ToolResponse, error) {
	request, err := buildGetK8sGetRequest(arguments)
	if err != nil {
		return nil, fmt.Errorf("error building k8s get request: %v", err)
	}

	body, err := getK8sGetMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting k8s resource snapshot: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func buildGetK8sGetRequest(arguments GetK8sGetHandlerArgs) (GetK8sGetRequest, error) {
	if err := validateRequiredString(arguments.ResourceAPIVersion, "resource_api_version"); err != nil {
		return GetK8sGetRequest{}, err
	}
	if err := validateRequiredString(arguments.ResourceKind, "resource_kind"); err != nil {
		return GetK8sGetRequest{}, err
	}
	if err := validateRequiredString(arguments.Name, "name"); err != nil {
		return GetK8sGetRequest{}, err
	}

	_, endTimeMs, err := calculateTimeRangeMillis(arguments.TimeConfig)
	if err != nil {
		return GetK8sGetRequest{}, fmt.Errorf("error calculating time range: %v", err)
	}

	format, err := normalizeK8sGetFormat(arguments.Format)
	if err != nil {
		return GetK8sGetRequest{}, err
	}

	request := GetK8sGetRequest{
		Environment: normalizeOptionalStringPtr(arguments.Environment),
		Namespace:   normalizeOptionalStringPtr(arguments.Namespace),
		Resource: k8sResourceReference{
			APIVersion: strings.TrimSpace(arguments.ResourceAPIVersion),
			Kind:       strings.TrimSpace(arguments.ResourceKind),
		},
		Name:   strings.TrimSpace(arguments.Name),
		UID:    normalizeOptionalStringPtr(arguments.UID),
		Time:   &endTimeMs,
		Format: format,
	}

	return request, nil
}

func normalizeK8sGetFormat(format string) (*string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(format))
	if trimmed == "" {
		return nil, nil
	}
	switch trimmed {
	case "yaml", "json":
		return &trimmed, nil
	default:
		return nil, fmt.Errorf("format must be yaml or json")
	}
}

func getK8sGetMetoroCall(ctx context.Context, request GetK8sGetRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling k8s get request: %v", err)
	}

	return utils.MakeMetoroAPIRequest("POST", "k8s/get", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
