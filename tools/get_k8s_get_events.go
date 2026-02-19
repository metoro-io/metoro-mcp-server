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

type GetK8sGetEventsHandlerArgs struct {
	TimeConfig         utils.TimeConfig `json:"time_config" jsonschema:"required,description=Time settings for this query. Start and end are used for event lookup"`
	Environment        string           `json:"environment" jsonschema:"description=Optional environment filter for this query"`
	Namespace          string           `json:"namespace" jsonschema:"description=Optional namespace filter for this query"`
	ResourceAPIVersion string           `json:"resource_api_version" jsonschema:"required,description=API version of the kubernetes resource such as v1 or apps v1"`
	ResourceKind       string           `json:"resource_kind" jsonschema:"required,description=Kind of the kubernetes resource such as Pod Deployment StatefulSet"`
	Name               string           `json:"name" jsonschema:"required,description=Resource name to fetch events for"`
	UID                string           `json:"uid" jsonschema:"description=Optional resource uid to disambiguate when names are reused"`
	Limit              *int             `json:"limit" jsonschema:"description=Optional page size. Must be greater than zero"`
	NextPageToken      string           `json:"next_page_token" jsonschema:"description=Optional token from previous page response"`
}

type GetK8sGetEventsRequest struct {
	Environment   *string              `json:"environment,omitempty"`
	Namespace     *string              `json:"namespace,omitempty"`
	Resource      k8sResourceReference `json:"resource"`
	Name          string               `json:"name"`
	UID           *string              `json:"uid,omitempty"`
	StartTime     int64                `json:"startTime"`
	EndTime       int64                `json:"endTime"`
	Limit         *int                 `json:"limit,omitempty"`
	NextPageToken *string              `json:"nextPageToken,omitempty"`
}

func GetK8sGetEventsHandler(ctx context.Context, arguments GetK8sGetEventsHandlerArgs) (*mcpgolang.ToolResponse, error) {
	request, err := buildGetK8sGetEventsRequest(arguments)
	if err != nil {
		return nil, fmt.Errorf("error building k8s get events request: %v", err)
	}

	body, err := getK8sGetEventsMetoroCall(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error getting k8s resource events: %v", err)
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(body)))), nil
}

func buildGetK8sGetEventsRequest(arguments GetK8sGetEventsHandlerArgs) (GetK8sGetEventsRequest, error) {
	if err := validateRequiredString(arguments.ResourceAPIVersion, "resource_api_version"); err != nil {
		return GetK8sGetEventsRequest{}, err
	}
	if err := validateRequiredString(arguments.ResourceKind, "resource_kind"); err != nil {
		return GetK8sGetEventsRequest{}, err
	}
	if err := validateRequiredString(arguments.Name, "name"); err != nil {
		return GetK8sGetEventsRequest{}, err
	}

	startTimeMs, endTimeMs, err := calculateTimeRangeMillis(arguments.TimeConfig)
	if err != nil {
		return GetK8sGetEventsRequest{}, fmt.Errorf("error calculating time range: %v", err)
	}

	request := GetK8sGetEventsRequest{
		Environment: normalizeOptionalStringPtr(arguments.Environment),
		Namespace:   normalizeOptionalStringPtr(arguments.Namespace),
		Resource: k8sResourceReference{
			APIVersion: strings.TrimSpace(arguments.ResourceAPIVersion),
			Kind:       strings.TrimSpace(arguments.ResourceKind),
		},
		Name:          strings.TrimSpace(arguments.Name),
		UID:           normalizeOptionalStringPtr(arguments.UID),
		StartTime:     startTimeMs,
		EndTime:       endTimeMs,
		Limit:         normalizeOptionalPositiveIntPtr(arguments.Limit),
		NextPageToken: normalizeOptionalStringPtr(arguments.NextPageToken),
	}

	return request, nil
}

func getK8sGetEventsMetoroCall(ctx context.Context, request GetK8sGetEventsRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling k8s get events request: %v", err)
	}

	return utils.MakeMetoroAPIRequest("POST", "k8s/getEvents", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
