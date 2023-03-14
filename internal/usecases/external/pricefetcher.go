package external

type PriceFetcher interface {
	GetPriceUSD(token string) (float32, error)
}
