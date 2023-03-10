package timeseries

import "errors"

func CalcSMAFromTimeSeries(series []*TimeValue) (float32, error) {
	var total float32 = 0.0
	for _, s := range series {
		total += s.Value
	}
	if len(series) > 0 {
		return total / (float32)(len(series)), nil
	}
	return 0, errors.New("series is too short")
}

func CalcEMAFromTimeSeries(series []*TimeValue, duration int, smooth float32) (float32, error) {
	if len(series) == 0 || len(series) < duration {
		return 0, errors.New("series is too short")
	}
	if len(series) == duration {
		return CalcSMAFromTimeSeries(series)
	}
	df := (float32)(duration)
	ema, _ := CalcSMAFromTimeSeries(series[:duration])
	for i := duration; i < len(series); i++ {
		ema = series[i].Value*smooth/(1+df) + ema*(1-smooth/(1+df))
	}

	return ema, nil
}
