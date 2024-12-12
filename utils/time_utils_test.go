package utils

import (
	"testing"
	"time"
)

func TestCalculateTimeRange(t *testing.T) {
	tests := []struct {
		name       string
		config     TimeConfig
		wantPeriod time.Duration // Expected time difference between start and end
	}{
		// Basic functionality tests
		{
			name: "5 minutes",
			config: TimeConfig{
				TimePeriod: 5,
				TimeWindow: "Minutes",
			},
			wantPeriod: 5 * time.Minute,
		},
		{
			name: "2 hours",
			config: TimeConfig{
				TimePeriod: 2,
				TimeWindow: "Hours",
			},
			wantPeriod: 2 * time.Hour,
		},
		{
			name: "3 days",
			config: TimeConfig{
				TimePeriod: 3,
				TimeWindow: "Days",
			},
			wantPeriod: 3 * 24 * time.Hour,
		},
		// Case sensitivity tests for Minutes
		{
			name: "minutes lowercase",
			config: TimeConfig{
				TimePeriod: 5,
				TimeWindow: "minutes",
			},
			wantPeriod: 5 * time.Minute,
		},
		{
			name: "minutes uppercase",
			config: TimeConfig{
				TimePeriod: 5,
				TimeWindow: "MINUTES",
			},
			wantPeriod: 5 * time.Minute,
		},
		{
			name: "minute singular",
			config: TimeConfig{
				TimePeriod: 1,
				TimeWindow: "minute",
			},
			wantPeriod: 1 * time.Minute,
		},
		{
			name: "mins abbreviation",
			config: TimeConfig{
				TimePeriod: 5,
				TimeWindow: "mins",
			},
			wantPeriod: 5 * time.Minute,
		},
		// Case sensitivity tests for Hours
		{
			name: "hours lowercase",
			config: TimeConfig{
				TimePeriod: 2,
				TimeWindow: "hours",
			},
			wantPeriod: 2 * time.Hour,
		},
		{
			name: "hours uppercase",
			config: TimeConfig{
				TimePeriod: 2,
				TimeWindow: "HOURS",
			},
			wantPeriod: 2 * time.Hour,
		},
		{
			name: "hour singular",
			config: TimeConfig{
				TimePeriod: 1,
				TimeWindow: "hour",
			},
			wantPeriod: 1 * time.Hour,
		},
		{
			name: "hrs abbreviation",
			config: TimeConfig{
				TimePeriod: 2,
				TimeWindow: "hrs",
			},
			wantPeriod: 2 * time.Hour,
		},
		// Case sensitivity tests for Days
		{
			name: "days lowercase",
			config: TimeConfig{
				TimePeriod: 3,
				TimeWindow: "days",
			},
			wantPeriod: 3 * 24 * time.Hour,
		},
		{
			name: "days uppercase",
			config: TimeConfig{
				TimePeriod: 3,
				TimeWindow: "DAYS",
			},
			wantPeriod: 3 * 24 * time.Hour,
		},
		{
			name: "day singular",
			config: TimeConfig{
				TimePeriod: 1,
				TimeWindow: "day",
			},
			wantPeriod: 1 * 24 * time.Hour,
		},
		{
			name: "default to minutes when window unspecified",
			config: TimeConfig{
				TimePeriod: 10,
				TimeWindow: "InvalidWindow",
			},
			wantPeriod: 10 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime, endTime := CalculateTimeRange(tt.config)
			
			// Convert Unix timestamps back to time.Time for easier comparison
			startTimeObj := time.Unix(startTime, 0)
			endTimeObj := time.Unix(endTime, 0)

			// Check if the time difference matches expected duration
			gotPeriod := endTimeObj.Sub(startTimeObj)
			if gotPeriod != tt.wantPeriod {
				t.Errorf("CalculateTimeRange() time period = %v, want %v", gotPeriod, tt.wantPeriod)
			}

			// Check if endTime is approximately now (within 1 second tolerance)
			nowUnix := time.Now().Unix()
			if diff := abs(endTime - nowUnix); diff > 1 {
				t.Errorf("CalculateTimeRange() endTime is not close enough to current time. diff = %v seconds", diff)
			}
		})
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
