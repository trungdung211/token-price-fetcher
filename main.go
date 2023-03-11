package main

import (
	"fmt"
	"time"

	"github.com/trungdung211/token-price-fetcher/pkg/timeseries"
)

// func main() {
// 	var username = "BotUser"
// 	var content = "☀ This is a test <strong>message</strong>"
// 	var url = "https://discord.com/api/webhooks/1083782248412753941/0S5RPsN0jR00cUa7ZlHL_zHzR01QTeqoZHqCJ-bv_on9jOSzOU6X26pKqsuSe8hvwoFH"

// 	message := discordwebhook.Message{
// 		Username: &username,
// 		Content:  &content,
// 	}

// 	err := discordwebhook.SendMessage(url, message)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	ts := time.Now()
	series := []*timeseries.TimeValue{
		{1, ts},
		{1, ts.Add(1 * time.Minute)},
		{2, ts.Add(2 * time.Minute)},
		{2, ts.Add(3 * time.Minute)},
		{3, ts.Add(5 * time.Minute)},
		{0.5, ts.Add(6 * time.Minute)},
	}

	ema, _ := timeseries.CalcEMAFromTimeSeries(series, 4, 2)
	fmt.Printf("%v\n", ema)
}
