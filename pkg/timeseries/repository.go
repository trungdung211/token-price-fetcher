package timeseries

type TimeSeriesRepository struct {
	timeseries           map[string]*MultiResolutionTimeSeries
	insertChan           *chan (*TimeValueResolution)
	capacity             int
	supportedResolutions []Resolution
}

func NewTimeSeriesRepository(insertChan *chan (*TimeValueResolution), capacity int, supportedResolutions []Resolution) *TimeSeriesRepository {
	return &TimeSeriesRepository{
		timeseries:           make(map[string]*MultiResolutionTimeSeries, 0),
		insertChan:           insertChan,
		capacity:             capacity,
		supportedResolutions: supportedResolutions,
	}
}

func (r *TimeSeriesRepository) Register(name string, series []*TimeValue) (err error) {
	mrt := NewMultiResolutionTimeSeries(name, r.supportedResolutions, r.capacity, r.insertChan)
	mrt.ReplaceSeries(series)

	r.timeseries[name] = mrt
	return
}
