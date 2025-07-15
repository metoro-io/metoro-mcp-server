package utils

import (
	"testing"
	"time"
	"strings"
)

func TestCalculateTimeRange(t *testing.T) {
	// Helper function to create pointers
	strPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }
	timeWindowPtr := func(tw TimeWindow) *TimeWindow { return &tw }

	tests := []struct {
		name          string
		config        TimeConfig
		wantPeriod    *time.Duration // Expected time difference between start and end for relative time
		wantStartTime *time.Time     // Expected exact start time for absolute time
		wantEndTime   *time.Time     // Expected exact end time for absolute time
		wantErr       bool
		errMsg        string
	}{
		// Relative time tests
		{
			name: "relative - 5 minutes",
			config: TimeConfig{
				Type:       RelativeTimeRange,
				TimePeriod: intPtr(5),
				TimeWindow: timeWindowPtr(Minutes),
			},
			wantPeriod: func() *time.Duration {
				d := 5 * time.Minute
				return &d
			}(),
		},
		{
			name: "relative - 2 hours",
			config: TimeConfig{
				Type:       RelativeTimeRange,
				TimePeriod: intPtr(2),
				TimeWindow: timeWindowPtr(Hours),
			},
			wantPeriod: func() *time.Duration {
				d := 2 * time.Hour
				return &d
			}(),
		},
		{
			name: "relative - missing time period",
			config: TimeConfig{
				Type:       RelativeTimeRange,
				TimeWindow: timeWindowPtr(Minutes),
			},
			wantErr: true,
			errMsg:  "time_period and time_window are required for relative time range",
		},
		{
			name: "relative - invalid time window",
			config: TimeConfig{
				Type:       RelativeTimeRange,
				TimePeriod: intPtr(5),
				TimeWindow: timeWindowPtr("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid time window: invalid",
		},
		{
			name: "relative - exceeds 30 days",
			config: TimeConfig{
				Type:       RelativeTimeRange,
				TimePeriod: intPtr(31),
				TimeWindow: timeWindowPtr(Days),
			},
			wantErr: true,
			errMsg:  "time range cannot exceed 30 days",
		},
		{
			name: "relative - exactly 30 days",
			config: TimeConfig{
				Type:       RelativeTimeRange,
				TimePeriod: intPtr(30),
				TimeWindow: timeWindowPtr(Days),
			},
			wantPeriod: func() *time.Duration {
				d := 30 * 24 * time.Hour
				return &d
			}(),
		},

		// Absolute time tests
		{
			name: "absolute - valid time range",
			config: TimeConfig{
				Type:      AbsoluteTimeRange,
				StartTime: strPtr(time.Now().Add(-24 * time.Hour).Format(time.RFC3339)),
				EndTime:   strPtr(time.Now().Add(-23 * time.Hour).Format(time.RFC3339)),
			},
			wantStartTime: func() *time.Time {
				t := time.Now().Add(-24 * time.Hour)
				return &t
			}(),
			wantEndTime: func() *time.Time {
				t := time.Now().Add(-23 * time.Hour)
				return &t
			}(),
		},
		{
			name: "absolute - missing start time",
			config: TimeConfig{
				Type:    AbsoluteTimeRange,
				EndTime: strPtr("2024-12-12T15:00:00Z"),
			},
			wantErr: true,
			errMsg:  "start_time and end_time are required for absolute time range",
		},
		{
			name: "absolute - invalid start time format",
			config: TimeConfig{
				Type:      AbsoluteTimeRange,
				StartTime: strPtr("invalid"),
				EndTime:   strPtr("2024-12-12T15:00:00Z"),
			},
			wantErr: true,
			errMsg:  "invalid start_time format",
		},
		{
			name: "absolute - end time before start time",
			config: TimeConfig{
				Type:      AbsoluteTimeRange,
				StartTime: strPtr("2024-12-12T15:00:00Z"),
				EndTime:   strPtr("2024-12-12T14:00:00Z"),
			},
			wantErr: true,
			errMsg:  "end_time cannot be before start_time",
		},
		{
			name: "absolute - exceeds 30 days",
			config: TimeConfig{
				Type:      AbsoluteTimeRange,
				StartTime: strPtr(time.Now().Add(-31 * 24 * time.Hour).Format(time.RFC3339)),
				EndTime:   strPtr(time.Now().Format(time.RFC3339)),
			},
			wantErr: true,
			errMsg:  "time range cannot exceed 30 days",
		},
		{
			name: "absolute - exactly 30 days ago",
			config: TimeConfig{
				Type:      AbsoluteTimeRange,
				StartTime: strPtr(time.Now().Add(-30*24*time.Hour + 1*time.Hour).Format(time.RFC3339)),
				EndTime:   strPtr(time.Now().Add(-29 * 24 * time.Hour).Format(time.RFC3339)),
			},
			wantStartTime: func() *time.Time {
				t := time.Now().Add(-30*24*time.Hour + 1*time.Hour)
				return &t
			}(),
			wantEndTime: func() *time.Time {
				t := time.Now().Add(-29 * 24 * time.Hour)
				return &t
			}(),
		},

		// Invalid type test
		{
			name: "invalid time range type",
			config: TimeConfig{
				Type: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid time range type: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime, endTime, err := CalculateTimeRange(tt.config)

			// Check error cases
			if tt.wantErr {
				if err == nil {
					t.Errorf("CalculateTimeRange() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("CalculateTimeRange() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("CalculateTimeRange() unexpected error = %v", err)
				return
			}

			// For relative time tests
			if tt.wantPeriod != nil {
				startTimeObj := time.Unix(startTime, 0)
				endTimeObj := time.Unix(endTime, 0)
				gotPeriod := endTimeObj.Sub(startTimeObj)
				if gotPeriod != *tt.wantPeriod {
					t.Errorf("CalculateTimeRange() time period = %v, want %v", gotPeriod, *tt.wantPeriod)
				}

				// Check if endTime is approximately now (within 1 second tolerance)
				nowUnix := time.Now().Unix()
				if diff := abs(endTime - nowUnix); diff > 1 {
					t.Errorf("CalculateTimeRange() endTime is not close enough to current time. diff = %v seconds", diff)
				}
			}

			// For absolute time tests
			if tt.wantStartTime != nil && tt.wantEndTime != nil {
				// For dynamic times (using time.Now()), check time range validity and approximate values
				if strings.Contains(tt.name, "30 days ago") || strings.Contains(tt.name, "valid time range") {
					// Just check the time range is valid
					if endTime < startTime {
						t.Errorf("CalculateTimeRange() endTime < startTime")
					}
					// Check approximate values (within 2 seconds tolerance)
					if diff := abs(startTime - tt.wantStartTime.Unix()); diff > 2 {
						t.Errorf("CalculateTimeRange() startTime difference too large: %v seconds", diff)
					}
					if diff := abs(endTime - tt.wantEndTime.Unix()); diff > 2 {
						t.Errorf("CalculateTimeRange() endTime difference too large: %v seconds", diff)
					}
				} else {
					// For static times, check exact match
					if startTime != tt.wantStartTime.Unix() {
						t.Errorf("CalculateTimeRange() startTime = %v, want %v", startTime, tt.wantStartTime.Unix())
					}
					if endTime != tt.wantEndTime.Unix() {
						t.Errorf("CalculateTimeRange() endTime = %v, want %v", endTime, tt.wantEndTime.Unix())
					}
				}
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
