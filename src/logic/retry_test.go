package logic_test

import (
	"errors"
	"testing"
	"time"

	"github.com/nag0yan/ideas/logic"
)

type CountTimer struct {
	Count int
}

func (c *CountTimer) Sleep(d time.Duration) {
	c.Count++
}

func TestRetryCount(t *testing.T) {
	tests := []struct {
		name           string
		fn             func() (int, error)
		maxRetries     int
		delay          time.Duration
		wantRetryCount int
	}{
		{
			name: "successful retry",
			fn: func() (int, error) {
				return 42, nil
			},
			maxRetries:     3,
			delay:          100 * time.Millisecond,
			wantRetryCount: 0,
		},
		{
			name: "failed retry",
			fn: func() (int, error) {
				return 0, errors.New("error")
			},
			maxRetries:     3,
			delay:          100 * time.Millisecond,
			wantRetryCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := &CountTimer{}
			logic.Retry(tt.fn, tt.
				maxRetries, tt.delay, timer)
			actualRetryCount := timer.Count
			if tt.wantRetryCount != actualRetryCount {
				t.Errorf("Retry() retry count = %v, want %v", actualRetryCount, tt.wantRetryCount)
			}
		})
	}
}

type DurationAccumulator struct {
	Total int
}

func (a *DurationAccumulator) Sleep(d time.Duration) {
	a.Total += int(d.Milliseconds())
}

func TestRetryTotalDuration(t *testing.T) {
	tests := []struct {
		name       string
		fn         func() (int, error)
		maxRetries int
		delay      time.Duration
		wantTotal  int
	}{
		{
			name: "successful retry",
			fn: func() (int, error) {
				return 42, nil
			},
			maxRetries: 3,
			delay:      100 * time.Millisecond,
			wantTotal:  0,
		},
		{
			name: "failed retry",
			fn: func() (int, error) {
				return 0, errors.New("error")
			},
			maxRetries: 3,
			delay:      100 * time.Millisecond,
			wantTotal:  100 + 200 + 400,
		},
		{
			name: "failed retry with large delay",
			fn: func() (int, error) {
				return 0, errors.New("error")
			},
			maxRetries: 3,
			delay:      50000 * time.Millisecond,
			wantTotal:  50000 + 100000 + 100000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accumulator := &DurationAccumulator{}
			logic.Retry(tt.fn, tt.maxRetries, tt.delay, accumulator)
			actualTotal := accumulator.Total
			if tt.wantTotal != actualTotal {
				t.Errorf("Retry() total duration = %v, want %v", actualTotal, tt.wantTotal)
			}
		})
	}
}
