package tools

import (
	"fmt"
	"strings"

	"github.com/metoro-io/metoro-mcp-server/utils"
)

const (
	k8sTimeModePoint = "point"
	k8sTimeModeRange = "range"
)

type k8sResourceReference struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

func calculateTimeRangeMillis(timeConfig utils.TimeConfig) (int64, int64, error) {
	startTimeSeconds, endTimeSeconds, err := utils.CalculateTimeRange(timeConfig)
	if err != nil {
		return 0, 0, err
	}
	return startTimeSeconds * 1000, endTimeSeconds * 1000, nil
}

func normalizeTimeMode(timeMode string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(timeMode))
	switch normalized {
	case k8sTimeModePoint:
		return k8sTimeModePoint, nil
	case k8sTimeModeRange:
		return k8sTimeModeRange, nil
	default:
		return "", fmt.Errorf("time_mode must be point or range")
	}
}

func normalizeOptionalStringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeOptionalPositiveIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	if *value <= 0 {
		return nil
	}
	return value
}

func validateRequiredString(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}
