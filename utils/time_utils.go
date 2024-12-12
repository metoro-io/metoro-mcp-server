package utils

import "time"

// TimeWindow represents supported time window units
type TimeWindow string

const (
	Minutes TimeWindow = "Minutes"
	Hours   TimeWindow = "Hours"
	Days    TimeWindow = "Days"
)

// TimeConfig holds the configuration for time range calculation
type TimeConfig struct {
	TimePeriod int        `json:"time_period" jsonschema:"required,description=The time period"`
	TimeWindow TimeWindow `json:"time_window" jsonschema:"required,description=The time window, e.g., Minutes, Hours, Days"`
}

// calculateTimeRange returns start and end timestamps based on the time configuration
func CalculateTimeRange(config TimeConfig) (startTime, endTime int64) {
	now := time.Now()
	var duration time.Duration

	switch config.TimeWindow {
	case Minutes:
		duration = time.Duration(config.TimePeriod) * time.Minute
	case Hours:
		duration = time.Duration(config.TimePeriod) * time.Hour
	case Days:
		duration = time.Duration(config.TimePeriod) * 24 * time.Hour
	default:
		// Default to minutes if unspecified
		duration = time.Duration(config.TimePeriod) * time.Minute
	}

	startTimeObj := now.Add(-duration)
	return startTimeObj.Unix(), now.Unix()
}
