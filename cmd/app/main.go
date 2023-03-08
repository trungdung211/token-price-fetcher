package main

import (
	"github.com/trungdung211/token-price-fetcher/internal/app"
	"github.com/trungdung211/token-price-fetcher/pkg/config"
)

func main() {
	// Configuration
	config.InitConfigs("", "")

	// Run
	app.Run()
}
