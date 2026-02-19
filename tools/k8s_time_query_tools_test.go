package tools

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TestBuildGetK8sListRequestPointUsesTime(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"
	limit := 25

	request, err := buildGetK8sListRequest(GetK8sListHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		TimeMode:           "point",
		Environment:        "prod",
		Namespace:          "default",
		ResourceAPIVersion: "v1",
		ResourceKind:       "Pod",
		Limit:              &limit,
		NextPageToken:      "abc123",
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if request.Time == nil {
		t.Fatalf("expected time to be set")
	}
	if request.StartTime != nil || request.EndTime != nil {
		t.Fatalf("expected startTime and endTime to be nil in point mode")
	}

	expectedEnd := mustParseRFC3339(t, end).UnixMilli()
	if *request.Time != expectedEnd {
		t.Fatalf("expected time %d got %d", expectedEnd, *request.Time)
	}

	if request.Resource.APIVersion != "v1" || request.Resource.Kind != "Pod" {
		t.Fatalf("unexpected resource object %+v", request.Resource)
	}
	if request.Limit == nil || *request.Limit != 25 {
		t.Fatalf("expected limit to be set")
	}
	if request.NextPageToken == nil || *request.NextPageToken != "abc123" {
		t.Fatalf("expected next page token to be set")
	}
}

func TestBuildGetK8sListRequestRangeUsesStartAndEnd(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	request, err := buildGetK8sListRequest(GetK8sListHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		TimeMode:           "range",
		ResourceAPIVersion: "apps/v1",
		ResourceKind:       "Deployment",
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if request.Time != nil {
		t.Fatalf("expected time to be nil in range mode")
	}
	if request.StartTime == nil || request.EndTime == nil {
		t.Fatalf("expected startTime and endTime to be set")
	}

	expectedStart := mustParseRFC3339(t, start).UnixMilli()
	expectedEnd := mustParseRFC3339(t, end).UnixMilli()
	if *request.StartTime != expectedStart {
		t.Fatalf("expected startTime %d got %d", expectedStart, *request.StartTime)
	}
	if *request.EndTime != expectedEnd {
		t.Fatalf("expected endTime %d got %d", expectedEnd, *request.EndTime)
	}
}

func TestBuildGetK8sListRequestSkipsPaginationWhenNotSet(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	request, err := buildGetK8sListRequest(GetK8sListHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		TimeMode:           "range",
		ResourceAPIVersion: "v1",
		ResourceKind:       "Pod",
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if request.Limit != nil {
		t.Fatalf("expected limit to be nil when not set")
	}
	if request.NextPageToken != nil {
		t.Fatalf("expected nextPageToken to be nil when not set")
	}
}

func TestBuildGetK8sGetRequestUsesEndTime(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	request, err := buildGetK8sGetRequest(GetK8sGetHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		ResourceAPIVersion: "v1",
		ResourceKind:       "Pod",
		Name:               "mypod",
		Format:             "json",
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	expectedEnd := mustParseRFC3339(t, end).UnixMilli()
	if request.Time == nil || *request.Time != expectedEnd {
		t.Fatalf("expected time to match end time")
	}
	if request.Format == nil || *request.Format != "json" {
		t.Fatalf("expected format json")
	}
}

func TestBuildGetK8sGetEventsRequestUsesRangeAndPagination(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"
	limit := 10

	request, err := buildGetK8sGetEventsRequest(GetK8sGetEventsHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type:      utils.AbsoluteTimeRange,
			StartTime: &start,
			EndTime:   &end,
		},
		ResourceAPIVersion: "v1",
		ResourceKind:       "Pod",
		Name:               "mypod",
		Limit:              &limit,
		NextPageToken:      "token-1",
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	expectedStart := mustParseRFC3339(t, start).UnixMilli()
	expectedEnd := mustParseRFC3339(t, end).UnixMilli()
	if request.StartTime != expectedStart {
		t.Fatalf("expected startTime %d got %d", expectedStart, request.StartTime)
	}
	if request.EndTime != expectedEnd {
		t.Fatalf("expected endTime %d got %d", expectedEnd, request.EndTime)
	}
	if request.Limit == nil || *request.Limit != 10 {
		t.Fatalf("expected limit to be set")
	}
	if request.NextPageToken == nil || *request.NextPageToken != "token-1" {
		t.Fatalf("expected nextPageToken to be set")
	}
}

func TestNewK8sTimeToolArgumentDescriptionsDoNotContainCommas(t *testing.T) {
	typesToCheck := []reflect.Type{
		reflect.TypeOf(GetK8sListHandlerArgs{}),
		reflect.TypeOf(GetK8sGetHandlerArgs{}),
		reflect.TypeOf(GetK8sGetEventsHandlerArgs{}),
	}

	for _, typ := range typesToCheck {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonSchemaTag := field.Tag.Get("jsonschema")
			description := extractDescription(jsonSchemaTag)
			if description == "" {
				continue
			}
			if strings.Contains(description, ",") {
				t.Fatalf("field %s in %s has comma in description: %q", field.Name, typ.Name(), description)
			}
		}
	}
}

func extractDescription(jsonSchemaTag string) string {
	parts := strings.SplitN(jsonSchemaTag, "description=", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func mustParseRFC3339(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}
	return parsed
}
