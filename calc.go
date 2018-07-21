package main

import (
	"fmt"
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
	metrics.Time.P50 = timeslices.percentile(0.5)
	metrics.Time.P75 = timeslices.percentile(0.75)
	metrics.Time.P95 = timeslices.percentile(0.95)
	metrics.Time.P99 = timeslices.percentile(0.99)
	metrics.Time.P999 = timeslices.percentile(0.999)
	metrics.Time.Max = timeslices.max()
	metrics.Time.Min = timeslices.min()
	c.Unlock()
	return metrics
}

func (ts timeSlice) Len() int {
	return len(ts)
}

func (ts timeSlice) Less(i, j int) bool {
	return int64(ts[i]) < int64(ts[j])
}

func (ts timeSlice) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
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
	return time.Duration(total.Nanoseconds() / int64(ts.Len()))
}

func (ts timeSlice) stdDev() time.Duration {
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

func (ts timeSlice) percentile(p float64) time.Duration {
	p = math.Min(float64(100), math.Max(float64(0), p))
	index := (p / float64(100)) * float64((ts.Len() - 1))
	fractionPart := index - math.Floor(index)
	intPart := int64(math.Floor(index))
	result := ts[intPart].Nanoseconds()
	if fractionPart > 0 {
		result += int64(fractionPart * float64((ts[intPart+1].Nanoseconds() - ts[intPart].Nanoseconds())))
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%dns", result))
	if err != nil {
		panic(err)
	}
	return duration
}
