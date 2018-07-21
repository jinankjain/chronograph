package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	ErrEmptyConfig = "Error empty configuration"
)

// Chronograph contains some basic data structures needed to store benchmark metadata
type Chronograph struct {
	sync.Mutex
	Config     *Config
	Count      uint64
	Timeslices timeSlice
}

type timeSlice []time.Duration

// Config contains some basic configuration for the experiment like No. of repetitions
type Config struct {
	Repetitions uint64
}

type Metrics struct {
	Time struct {
		Cumulative time.Duration
		Mean       time.Duration
		P50        time.Duration
		P75        time.Duration
		P95        time.Duration
		P99        time.Duration
		P999       time.Duration
		Max        time.Duration
		Min        time.Duration
		StdDev     time.Duration
	}
	Sample uint64
	Count  uint64
}

func (m *Metrics) String() string {
	return fmt.Sprintf(`%d samples of %d events
Cumulative: %s
Mean:       %s
P50:        %s
P75:        %s
P95:        %s
P99:        %s
P999:       %s
Max:        %s
Min:        %s
StdDev:     %s`,
		m.Sample,
		m.Count,
		m.Time.Cumulative,
		m.Time.Mean,
		m.Time.P50,
		m.Time.P75,
		m.Time.P95,
		m.Time.P99,
		m.Time.P999,
		m.Time.Max,
		m.Time.Min,
		m.Time.StdDev)
}

// AddEvent add a duration to Chronograph
func (c *Chronograph) AddEvent(t time.Duration) {
	c.Lock()
	c.Timeslices = append(c.Timeslices, t)
	c.Count = c.Count + 1
	c.Unlock()
}

func NewChronograph(c *Config) (*Chronograph, error) {
	if c == nil {
		return nil, errors.New(ErrEmptyConfig)
	}
	return &Chronograph{Config: c, Count: 0}, nil
}
