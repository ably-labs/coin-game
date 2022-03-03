package handlers

import (
	// "context"
	"log"
	// "stocker/app/stocks"
	"stocker/config"

	"github.com/ably/ably-go/ably"
)

func AblyClient() *ably.Realtime {
	client, err := ably.NewRealtime(ably.WithKey(config.EnvVariable("API_KEY")))
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

// func RegisterPresense(ctx context.Context, chn *ably.RealtimeChannel, playerName string) {
// 	err := chn.Presence.Enter(ctx, playerName)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func Publish(ctx context.Context, chn *ably.RealtimeChannel, data *stocks.CoinResponse)
