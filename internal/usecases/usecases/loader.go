package usecases

import "github.com/trungdung211/token-price-fetcher/pkg/timeseries"

type historyLoader struct {
}

type HistoryLoader interface {
	Load(process func(name string, series []*timeseries.TimeValue) error, capacity int) error
	LoadOne(process func(name string, series []*timeseries.TimeValue) error, name string, capacity int) error
}

func NewHistoryLoader() HistoryLoader {
	return &historyLoader{}
}

func (h *historyLoader) Load(process func(name string, series []*timeseries.TimeValue) error, capacity int) error {
	return nil
}

func (h *historyLoader) LoadOne(process func(name string, series []*timeseries.TimeValue) error, name string, capacity int) error {
	return nil
}
