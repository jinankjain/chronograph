package main

import (
	"math"
	"sort"
	"time"
)

func (c *Chronograph) Calculate() *Metrics {
	metrics := &Metrics{}
	c.Lock()
	if c.Count == 0 {
		return metrics
	}
	metrics.Sample = c.Count
	metrics.Count = c.Config.Repetitions
	timeslices := make(timeSlice, metrics.Sample)
	copy(timeslices, c.Timeslices)
	sort.Sort(timeslices)
	metrics.Time.Cumulative = timeslices.cumulative()
	metrics.Time.Mean = timeslices.mean()
	metrics.Time.StdDev = timeslices.stdDev()
	metrics.Time.P50 = timeslices[timeslices.Len()/2]
	metrics.Time.P75 = timeslices.p(0.75)
	metrics.Time.P95 = timeslices.p(0.95)
	metrics.Time.P99 = timeslices.p(0.99)
	metrics.Time.P999 = timeslices.p(0.999)
	metrics.Time.Max = timeslices.max()
	metrics.Time.Min = timeslices.min()
	c.Unlock()
}

func (ts timeSlice) Len() int {
	return len(ts)
}

func (ts timeSlice) cumulative() time.Duration {
	var total time.Duration
	for _, t := range ts {
		total += t
	}
	return total
}

func (ts timeSlice) mean() time.Duration {
	var total time.Duration
	for _, t := range ts {
		total += t
	}
	return time.Duration(total.Nanoseconds() / ts.Len())
}

func (ts timeSlice) stdDev() time.Duration {
	var total time.Duration
	avg := ts.mean()
	total := 0.00
	for _, t := range ts {
		total += math.Pow(float64(avg-t), 2)
	}
	totalq := total / (float64(ts.Len()))
	return time.Duration(math.Sqrt(totalq))
}

func (ts timeSlice) min() time.Duration {
	return ts[0]
}

func (ts timeSlice) max() time.Duration {
	return ts[ts.Len()-1]
}

func (ts timeSlice) p(p float64) time.Duration {
	return ts[int(float64(ts.Len())*p+0.5)-1]
}
