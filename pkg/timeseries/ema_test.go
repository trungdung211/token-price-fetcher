package timeseries

import (
	"math"
	"testing"
	"time"
)

func TestEmaCalc(t *testing.T) {
	ts := time.Now()
	series := []*TimeValue{
		{1, ts},
		{1, ts.Add(1 * time.Minute)},
		{2, ts.Add(2 * time.Minute)},
		{2, ts.Add(3 * time.Minute)},
		{3, ts.Add(5 * time.Minute)},
		{0.5, ts.Add(6 * time.Minute)},
	}

	ema, _ := CalcEMAFromTimeSeries(series, 4, 2)
	expected := 1.4600002
	if math.Abs((float64)(ema)-expected) > 1e-7 {
		t.Errorf("Expect ema is %v instead of %v", expected, ema)
	}
}
