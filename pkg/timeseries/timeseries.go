package timeseries

import (
	"errors"
	"time"
)

type TimeValue struct {
	Value float32
	Time  time.Time
}

type TimeValueResolution struct {
	TV         *TimeValue
	Resolution Resolution
	Name       string
}

type MultiResolutionTimeSeries struct {
	name       string
	buckets    map[Resolution]*Bucket
	series     map[Resolution][]*TimeValue
	capacity   int
	insertChan *chan (*TimeValueResolution)
}

func NewMultiResolutionTimeSeries(name string, resolutions []Resolution, capacity int, insertChan *chan (*TimeValueResolution)) *MultiResolutionTimeSeries {
	buckets := make(map[Resolution]*Bucket, len(resolutions))
	series := make(map[Resolution][]*TimeValue, len(resolutions))
	for _, resol := range resolutions {
		buckets[resol] = NewBucket(resol)
	}
	return &MultiResolutionTimeSeries{
		name:       name,
		buckets:    buckets,
		series:     series,
		capacity:   capacity,
		insertChan: insertChan,
	}
}

func (t *MultiResolutionTimeSeries) Add(value float32, ts time.Time, new bool) {
	for resol, bucket := range t.buckets {
		if out, outTs, finished := bucket.Add(value, ts); finished {
			tv := &TimeValue{
				Value: out,
				Time:  outTs,
			}
			t.series[resol] = append(t.series[resol], tv)
			if len(t.series[resol]) > t.capacity {
				t.series[resol] = t.series[resol][1:]
			}
			if new && t.insertChan != nil {
				(*t.insertChan) <- &TimeValueResolution{
					TV:         tv,
					Resolution: resol,
					Name:       t.name,
				}
			}
		}
	}
}

func (t *MultiResolutionTimeSeries) ReplaceSeries(series []*TimeValue) (err error) {
	for _, s := range series {
		t.Add(s.Value, s.Time, false)
	}
	return
}

func (t *MultiResolutionTimeSeries) GetSeries(resolution Resolution) (out []*TimeValue, err error) {
	if series, found := t.series[resolution]; found {
		out = series
		return
	}

	err = errors.New("not found resolution")
	return
}

func (t *MultiResolutionTimeSeries) GetResolutions() []Resolution {
	keys := make([]Resolution, 0)
	for k, _ := range t.buckets {
		keys = append(keys, k)
	}

	return keys
}
