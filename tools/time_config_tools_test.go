package tools

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/metoro-io/metoro-mcp-server/utils"
)

func TestBuildGetServicesPathUsesTimeConfig(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	path, err := buildGetServicesPath(GetServicesHandlerArgs{
		TimeConfig: absoluteTimeConfig(start, end),
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	assertPathAndTimes(t, path, "services", start, end)
}

func TestBuildGetServicesPathReturnsErrorForInvalidTimeConfig(t *testing.T) {
	_, err := buildGetServicesPath(GetServicesHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type: utils.AbsoluteTimeRange,
		},
	})
	if err == nil {
		t.Fatalf("expected error for invalid time config")
	}
}

func TestBuildGetNamespacesPathUsesTimeConfig(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	path, err := buildGetNamespacesPath(GetNamespacesHandlerArgs{
		TimeConfig: absoluteTimeConfig(start, end),
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	assertPathAndTimes(t, path, "namespaces", start, end)
}

func TestBuildGetNamespacesPathReturnsErrorForInvalidTimeConfig(t *testing.T) {
	_, err := buildGetNamespacesPath(GetNamespacesHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type: utils.AbsoluteTimeRange,
		},
	})
	if err == nil {
		t.Fatalf("expected error for invalid time config")
	}
}

func TestBuildGetK8sEventsSummaryAttributesPathUsesTimeConfig(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	path, err := buildGetK8sEventsSummaryAttributesPath(GetK8sEventsAttributesHandlerArgs{
		TimeConfig: absoluteTimeConfig(start, end),
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	assertPathAndTimes(t, path, "k8s/events/summaryAttributes", start, end)
}

func TestBuildGetK8sEventsSummaryAttributesPathReturnsErrorForInvalidTimeConfig(t *testing.T) {
	_, err := buildGetK8sEventsSummaryAttributesPath(GetK8sEventsAttributesHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type: utils.AbsoluteTimeRange,
		},
	})
	if err == nil {
		t.Fatalf("expected error for invalid time config")
	}
}

func TestBuildGetNodeInfoPathUsesEndTimeAsPoint(t *testing.T) {
	start := "2026-02-19T10:00:00Z"
	end := "2026-02-19T10:05:00Z"

	path, err := buildGetNodeInfoPath(GetNodeInfoHandlerArgs{
		TimeConfig: absoluteTimeConfig(start, end),
		NodeName:   "node-1",
	})
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	parsed, err := url.Parse("http://example/" + path)
	if err != nil {
		t.Fatalf("failed to parse path: %v", err)
	}

	if parsed.Path != "/infrastructure/node" {
		t.Fatalf("unexpected path %q", parsed.Path)
	}

	expectedEnd := strconv.FormatInt(mustParseTime(t, end).Unix(), 10)
	if got := parsed.Query().Get("time"); got != expectedEnd {
		t.Fatalf("expected time=%s got %s", expectedEnd, got)
	}
	if got := parsed.Query().Get("startTime"); got != expectedEnd {
		t.Fatalf("expected startTime=%s got %s", expectedEnd, got)
	}
	if got := parsed.Query().Get("nodeName"); got != "node-1" {
		t.Fatalf("expected nodeName=node-1 got %s", got)
	}
}

func TestBuildGetNodeInfoPathReturnsErrorForInvalidTimeConfig(t *testing.T) {
	_, err := buildGetNodeInfoPath(GetNodeInfoHandlerArgs{
		TimeConfig: utils.TimeConfig{
			Type: utils.AbsoluteTimeRange,
		},
		NodeName: "node-1",
	})
	if err == nil {
		t.Fatalf("expected error for invalid time config")
	}
}

func assertPathAndTimes(t *testing.T, path, expectedPath, start, end string) {
	t.Helper()

	parsed, err := url.Parse("http://example/" + path)
	if err != nil {
		t.Fatalf("failed to parse path: %v", err)
	}

	if parsed.Path != "/"+expectedPath {
		t.Fatalf("unexpected path %q", parsed.Path)
	}

	expectedStart := strconv.FormatInt(mustParseTime(t, start).Unix(), 10)
	expectedEnd := strconv.FormatInt(mustParseTime(t, end).Unix(), 10)
	if got := parsed.Query().Get("startTime"); got != expectedStart {
		t.Fatalf("expected startTime=%s got %s", expectedStart, got)
	}
	if got := parsed.Query().Get("endTime"); got != expectedEnd {
		t.Fatalf("expected endTime=%s got %s", expectedEnd, got)
	}
}

func absoluteTimeConfig(start, end string) utils.TimeConfig {
	return utils.TimeConfig{
		Type:      utils.AbsoluteTimeRange,
		StartTime: &start,
		EndTime:   &end,
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", value, err)
	}
	return parsed
}
