package timeseries

import (
	"time"
)

type Resolution int

const (
	TIME_RESOLUTION_1_MIN Resolution = iota
	TIME_RESOLUTION_1_HOUR
	TIME_RESOLUTION_4_HOURS
	TIME_RESOLUTION_1_DAY
)

func (r *Resolution) ToString() string {
	switch *r {
	case TIME_RESOLUTION_1_MIN:
		return "1min"
	case TIME_RESOLUTION_1_HOUR:
		return "1hour"
	case TIME_RESOLUTION_4_HOURS:
		return "4hours"
	case TIME_RESOLUTION_1_DAY:
		return "1day"
	}
	return ""
}

func resolution2Duration(resol Resolution) time.Duration {
	switch resol {
	case TIME_RESOLUTION_1_MIN:
		return 1 * time.Minute
	case TIME_RESOLUTION_1_HOUR:
		return 1 * time.Hour
	case TIME_RESOLUTION_4_HOURS:
		return 4 * time.Hour
	default:
		return 24 * time.Hour
	}
}

func getTimeInResolution(t time.Time, resolutionDuration time.Duration) time.Time {
	ts := t.Unix()
	return t.Add(-time.Duration((ts % int64(resolutionDuration.Seconds()))) * time.Second)
}

type Bucket struct {
	resolution         Resolution
	resolutionDuration time.Duration
	nextTime           *time.Time
}

func NewBucket(resolution Resolution) *Bucket {
	return &Bucket{
		resolution:         resolution,
		resolutionDuration: resolution2Duration(resolution),
	}
}

func (b *Bucket) Add(value float32, timestamp time.Time) (out float32, t time.Time, finished bool) {
	if b.nextTime == nil || b.nextTime.Before(timestamp) {
		out = value
		t = getTimeInResolution(timestamp, b.resolutionDuration)
		nextTime := t.Add(b.resolutionDuration)
		b.nextTime = &nextTime
		finished = true
	}
	return
}
