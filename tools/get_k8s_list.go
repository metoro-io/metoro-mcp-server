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

type GetK8sListHandlerArgs struct {
	TimeConfig         utils.TimeConfig `json:"time_config" jsonschema:"required,description=Time settings for this query. Use relative or absolute window"`
	TimeMode           string           `json:"time_mode" jsonschema:"required,enum=point,enum=range,description=Choose point to query one timestamp or range to query overlap window"`
	Environment        string           `json:"environment" jsonschema:"description=Optional environment filter for this query"`
	Namespace          string           `json:"namespace" jsonschema:"description=Optional namespace filter for this query"`
	ResourceAPIVersion string           `json:"resource_api_version" jsonschema:"required,description=API version of the kubernetes resource such as v1 or apps v1"`
	ResourceKind       string           `json:"resource_kind" jsonschema:"required,description=Kind of the kubernetes resource such as Pod Deployment StatefulSet"`
	Limit              *int             `json:"limit" jsonschema:"description=Optional page size. Must be greater than zero"`
	NextPageToken      string           `json:"next_page_token" jsonschema:"description=Optional token from previous page response"`
}

type GetK8sListRequest struct {
	Environment   *string              `json:"environment,omitempty"`
	Namespace     *string              `json:"namespace,omitempty"`
	Resource      k8sResourceReference `json:"resource"`
	Time          *int64               `json:"time,omitempty"`
	StartTime     *int64               `json:"startTime,omitempty"`
	EndTime       *int64               `json:"endTime,omitempty"`
	Limit         *int                 `json:"limit,omitempty"`
	NextPageToken *string              `json:"nextPageToken,omitempty"`
}

func GetK8sListHandler(ctx context.Context, arguments GetK8sListHandlerArgs) (*mcpgolang.ToolResponse, error) {
	request, err := buildGetK8sListRequest(arguments)
	if err != nil {
		return nil, fmt.Errorf("error building k8s list request: %v", err)
	}

	body, err := getK8sListMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting k8s list: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func buildGetK8sListRequest(arguments GetK8sListHandlerArgs) (GetK8sListRequest, error) {
	if err := validateRequiredString(arguments.ResourceAPIVersion, "resource_api_version"); err != nil {
		return GetK8sListRequest{}, err
	}
	if err := validateRequiredString(arguments.ResourceKind, "resource_kind"); err != nil {
		return GetK8sListRequest{}, err
	}

	timeMode, err := normalizeTimeMode(arguments.TimeMode)
	if err != nil {
		return GetK8sListRequest{}, err
	}

	startTimeMs, endTimeMs, err := calculateTimeRangeMillis(arguments.TimeConfig)
	if err != nil {
		return GetK8sListRequest{}, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetK8sListRequest{
		Environment: normalizeOptionalStringPtr(arguments.Environment),
		Namespace:   normalizeOptionalStringPtr(arguments.Namespace),
		Resource: k8sResourceReference{
			APIVersion: strings.TrimSpace(arguments.ResourceAPIVersion),
			Kind:       strings.TrimSpace(arguments.ResourceKind),
		},
		Limit:         normalizeOptionalPositiveIntPtr(arguments.Limit),
		NextPageToken: normalizeOptionalStringPtr(arguments.NextPageToken),
	}

	if timeMode == k8sTimeModePoint {
		request.Time = &endTimeMs
	} else {
		request.StartTime = &startTimeMs
		request.EndTime = &endTimeMs
	}

	return request, nil
}

func getK8sListMetoroCall(ctx context.Context, request GetK8sListRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling k8s list request: %v", err)
	}

	return utils.MakeMetoroAPIRequest("POST", "k8s/list", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
