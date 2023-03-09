package timeseries

import (
	"testing"
	"time"
)

func TestMultiResolutions(t *testing.T) {
	mrt := NewMultiResolutionTimeSeries("", []Resolution{
		TIME_RESOLUTION_1_DAY,
		TIME_RESOLUTION_1_HOUR,
	}, 10, nil)

	if len(mrt.GetResolutions()) != 2 {
		t.Errorf("Expect number of resolutions is %v instead of %v", 2, len(mrt.GetResolutions()))
	}
}

func TestCapacity(t *testing.T) {
	capacity := 2

	mrt := NewMultiResolutionTimeSeries("", []Resolution{
		TIME_RESOLUTION_1_MIN,
	}, capacity, nil)

	ts := time.Now()
	mrt.Add(1.0, ts)
	mrt.Add(2.0, ts.Add(1*time.Minute))
	mrt.Add(3.0, ts.Add(2*time.Minute))

	series, _ := mrt.GetSeries(TIME_RESOLUTION_1_MIN)

	if len(series) != capacity {
		t.Errorf("Expect length of series is capacity %v instead of %v", capacity, len(series))
	}

	if series[len(series)-1].Value != 3.0 {
		t.Errorf("Expect series slice from right, last value should be %v instead of %v", 3.0, series[len(series)-1].Value)
	}
}

func TestBucket(t *testing.T) {
	capacity := 4

	mrt := NewMultiResolutionTimeSeries("", []Resolution{
		TIME_RESOLUTION_1_MIN,
		TIME_RESOLUTION_4_HOURS,
	}, capacity, nil)

	ts_base, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 16:00:00")
	ts := ts_base.Add(80 * time.Second)

	mrt.Add(1.0, ts)
	mrt.Add(2.0, ts.Add(1*time.Minute))
	mrt.Add(3.0, ts.Add(2*time.Minute))
	mrt.Add(4.0, ts.Add(2*time.Hour))
	mrt.Add(5.0, ts.Add(4*time.Hour))

	series_1_min, _ := mrt.GetSeries(TIME_RESOLUTION_1_MIN)

	if len(series_1_min) != capacity {
		t.Errorf("Expect length of series is capacity %v instead of %v", capacity, len(series_1_min))
	}

	if series_1_min[len(series_1_min)-1].Value != 5.0 {
		t.Errorf("Expect series slice from right, last value should be %v instead of %v", 5.0, series_1_min[len(series_1_min)-1].Value)
	}

	series_4_hours, _ := mrt.GetSeries(TIME_RESOLUTION_4_HOURS)

	if len(series_4_hours) != 2 {
		t.Errorf("Expect length of series is capacity %v instead of %v", 2, len(series_4_hours))
	}

	if series_4_hours[len(series_4_hours)-1].Value != 5.0 || series_4_hours[len(series_4_hours)-2].Value != 1.0 {
		t.Errorf("Expect series slice from right, 2 last values should be %v, %v instead of %v, %v", 1, 5, series_4_hours[len(series_4_hours)-2].Value, series_4_hours[len(series_4_hours)-1].Value)
	}
	if !series_4_hours[len(series_4_hours)-2].Time.Equal(ts_base) {
		t.Errorf("Expect time bucket is %v instead of %v", ts_base, series_4_hours[len(series_4_hours)-2].Time)
	}
}

func TestInitSeries(t *testing.T) {
	mrt := NewMultiResolutionTimeSeries("", []Resolution{
		TIME_RESOLUTION_1_MIN,
		TIME_RESOLUTION_1_HOUR,
		TIME_RESOLUTION_4_HOURS,
	}, 10, nil)

	ts_base, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 16:00:00")
	ts := ts_base.Add(80 * time.Second)

	mrt.ReplaceSeries([]*TimeValue{
		{Value: 1, Time: ts},
		{Value: 2, Time: ts.Add(1 * time.Minute)},
		{Value: 3, Time: ts.Add(1 * time.Hour)},
		{Value: 4, Time: ts.Add(2 * time.Hour)},
		{Value: 5, Time: ts.Add(3 * time.Hour)},
	})

	series_1_min, _ := mrt.GetSeries(TIME_RESOLUTION_1_MIN)
	series_1_hour, _ := mrt.GetSeries(TIME_RESOLUTION_1_HOUR)
	series_4_hours, _ := mrt.GetSeries(TIME_RESOLUTION_4_HOURS)

	if len(series_1_min) != 5 {
		t.Errorf("Expect length of series_1_min is capacity %v instead of %v", 5, len(series_1_min))
	}
	if len(series_1_hour) != 4 {
		t.Errorf("Expect length of series_1_hour is capacity %v instead of %v", 4, len(series_1_hour))
	}
	if len(series_4_hours) != 1 {
		t.Errorf("Expect length of series_4_hours is capacity %v instead of %v", 1, len(series_4_hours))
	}

}
