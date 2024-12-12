package utils

import (
	"fmt"
	"strings"
	"time"
)

// TimeWindow represents supported time window units
type TimeWindow string

const (
	Minutes TimeWindow = "Minutes"
	Hours   TimeWindow = "Hours"
	Days    TimeWindow = "Days"
)

// TimeRangeType indicates whether the time range is relative or absolute
type TimeRangeType string

const (
	RelativeTimeRange TimeRangeType = "relative"
	AbsoluteTimeRange TimeRangeType = "absolute"
)

// TimeConfig holds the configuration for time range calculation
type TimeConfig struct {
	// Type of time range (relative or absolute)
	Type TimeRangeType `json:"type" jsonschema:"required,enum=relative,enum=absolute,description=Type of time range. Must be either 'relative' or 'absolute'"`

	// Fields for relative time range
	TimePeriod *int       `json:"time_period,omitempty" jsonschema:"description=For relative time range: the number of time units to look back"`
	TimeWindow *TimeWindow `json:"time_window,omitempty" jsonschema:"description=For relative time range: the unit of time (Minutes, Hours, Days)"`

	// Fields for absolute time range
	StartTime *string `json:"start_time,omitempty" jsonschema:"description=For absolute time range: start time in RFC3339 format (e.g., '2024-12-12T14:27:22Z')"`
	EndTime   *string `json:"end_time,omitempty" jsonschema:"description=For absolute time range: end time in RFC3339 format (e.g., '2024-12-12T14:27:22Z')"`
}

// CalculateTimeRange returns start and end timestamps based on the time configuration
func CalculateTimeRange(config TimeConfig) (startTime, endTime int64, err error) {
	now := time.Now()

	switch config.Type {
	case RelativeTimeRange:
		if config.TimePeriod == nil || config.TimeWindow == nil {
			return 0, 0, fmt.Errorf("time_period and time_window are required for relative time range")
		}

		var duration time.Duration
		window := strings.ToLower(string(*config.TimeWindow))

		switch window {
		case "minutes", "minute", "min", "mins":
			duration = time.Duration(*config.TimePeriod) * time.Minute
		case "hours", "hour", "hr", "hrs":
			duration = time.Duration(*config.TimePeriod) * time.Hour
		case "days", "day":
			duration = time.Duration(*config.TimePeriod) * 24 * time.Hour
		default:
			return 0, 0, fmt.Errorf("invalid time window: %s", *config.TimeWindow)
		}

		startTimeObj := now.Add(-duration)
		return startTimeObj.Unix(), now.Unix(), nil

	case AbsoluteTimeRange:
		if config.StartTime == nil || config.EndTime == nil {
			return 0, 0, fmt.Errorf("start_time and end_time are required for absolute time range")
		}

		startTimeObj, err := time.Parse(time.RFC3339, *config.StartTime)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid start_time format: %v", err)
		}

		endTimeObj, err := time.Parse(time.RFC3339, *config.EndTime)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid end_time format: %v", err)
		}

		if endTimeObj.Before(startTimeObj) {
			return 0, 0, fmt.Errorf("end_time cannot be before start_time")
		}

		return startTimeObj.Unix(), endTimeObj.Unix(), nil

	default:
		return 0, 0, fmt.Errorf("invalid time range type: %s", config.Type)
	}
}
