package chronograph

import (
	"errors"
	"testing"
	"time"
)

func TestNewChronographApi(t *testing.T) {
	tests := []struct {
		config *Config
		err    error
	}{
		{config: nil, err: errors.New(ErrEmptyConfig)},
		{config: &Config{Repetitions: 30}, err: nil},
	}
	for _, test := range tests {
		c, err := NewChronograph(test.config)
		if err != nil {
			if test.err.Error() != err.Error() {
				t.Fatalf("Error did not match. Expected: %s, Actual: %s", test.err, err)
			}
		}
		if test.config != nil && test.config.Repetitions != c.Config.Repetitions {
			t.Fatalf("Error while creating new Chronograph. Expected: %d, Actual: %d",
				test.config.Repetitions, c.Config.Repetitions)
		}
	}
}

func doWork(n int) {
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func TestChronographAddEventApi(t *testing.T) {
	config := &Config{Repetitions: 40}
	chrono, err := NewChronograph(config)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(config.Repetitions); i++ {
		chrono.AddEvent(time.Millisecond)
	}
	metrics := chrono.Calculate()
	if chrono.Timeslices.Len() != int(config.Repetitions) {
		t.Fatalf("Error missed number of repetitions. Expected: %d, Actual: %d", config.Repetitions,
			chrono.Timeslices.Len())
	}
	if metrics.Time.Cumulative != (time.Millisecond * 40) {
		t.Fatalf("Mismatch cumulative sum. Expected: %d, Actual: %d", 40*time.Millisecond,
			metrics.Time.Cumulative)
	}
	if metrics.Time.StdDev != 0 {
		t.Fatalf("Mismatch stddev. Expected: %d, Actual: %d", 0, metrics.Time.StdDev)
	}
	if metrics.Time.Max != time.Millisecond {
		t.Fatalf("Mismatch max value. Expected: %d, Actual: %d", time.Microsecond,
			metrics.Time.Max)
	}
	if metrics.Time.Min != time.Millisecond {
		t.Fatalf("Mismatch min value. Expected: %d, Actual: %d", time.Microsecond,
			metrics.Time.Min)
	}
	if metrics.Time.P50 != time.Millisecond {
		t.Fatalf("Mismatch P50 value. Expected: %d, Actual: %d", time.Microsecond,
			metrics.Time.P50)
	}
	if metrics.Time.P75 != time.Millisecond {
		t.Fatalf("Mismatch P75 value. Expected: %d, Actual: %d", time.Microsecond,
			metrics.Time.P75)
	}
	if metrics.Time.P95 != time.Millisecond {
		t.Fatalf("Mismatch P95 value. Expected: %d, Actual: %d", time.Microsecond,
			metrics.Time.P95)
	}
	if metrics.Time.P99 != time.Millisecond {
		t.Fatalf("Mismatch P99 value. Expected: %d, Actual: %d", time.Microsecond,
			metrics.Time.P99)
	}
	if metrics.Time.P999 != time.Millisecond {
		t.Fatalf("Mismatch P999 value. Expected: %d, Actual: %d", time.Microsecond,
			metrics.Time.P999)
	}
}
