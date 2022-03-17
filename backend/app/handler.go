package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"stocker/config"

	"time"

	"github.com/ably/ably-go/ably"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

var (
	channel    = AblyClient().Channels.Get("stock")
	ctx        = context.Background()
	cacheStore = cache.New(5*time.Minute, 10*time.Minute)

	// WarningLogger ...
	WarningLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	// InfoLogger ...
	InfoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// ErrorLogger ...
	ErrorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Println("Creating player...")
	player := Player{}
	err := json.NewDecoder(r.Body).Decode(&player)
	player.Wallet = 700000
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheStore.Set(player.Name, &player, cache.DefaultExpiration)
	play, found := cacheStore.Get(player.Name)
	if found {
		player = play.(Player)
	}
	res := &CoinResponse{
		Player:        player.Name,
		CoinQuantity:  player.CoinQuantity,
		WalletBalance: player.Wallet,
	}

	writeResponse(res, w, r)
}

// GetWalletBalance - get current wallet balance
func GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Println("Getting wallet balance...")
	currentPlayer := Player{}
	var res *CoinResponse
	player := mux.Vars(r)["player"]

	play, found := cacheStore.Get(player)
	if found {
		currentPlayer = play.(Player)
	}
	res = &CoinResponse{
		Player:        player,
		CoinQuantity:  currentPlayer.CoinQuantity,
		WalletBalance: currentPlayer.Wallet,
	}
	writeResponse(res, w, r)
}

func BuyBitcoin(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Println("Buying coin...")
	currentPlayer := Player{}
	buyRequest := BuySellRequest{}
	err := json.NewDecoder(r.Body).Decode(&buyRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	play, found := cacheStore.Get(buyRequest.Player)
	if found {
		currentPlayer = play.(Player)
	}

	playerBalance := currentPlayer.Wallet
	playerQuantity := currentPlayer.CoinQuantity

	totalCost := buyRequest.CurrentCoinPrice * buyRequest.Quantity
	if totalCost > playerBalance {
		http.Error(w, "Insufficient wallet balance", http.StatusBadRequest)
		return
	} else {
		currentPlayer.Wallet = playerBalance - totalCost
		currentPlayer.CoinQuantity = playerQuantity + buyRequest.Quantity

		cacheStore.Set(currentPlayer.Name, currentPlayer, cache.DefaultExpiration)

		res := &CoinResponse{
			Player:        buyRequest.Player,
			CoinQuantity:  currentPlayer.CoinQuantity,
			WalletBalance: currentPlayer.Wallet,
		}

		channel.Publish(ctx, "buy", res)
		// writeResponse(res, w, r)
	}

}

func SellBitcoin(w http.ResponseWriter, r *http.Request) {
	InfoLogger.Println("Selling coin...")
	sellRequest := BuySellRequest{}
	currentPlayer := Player{}
	err := json.NewDecoder(r.Body).Decode(&sellRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	play, found := cacheStore.Get(sellRequest.Player)
	if found {
		currentPlayer = play.(Player)
	}

	playerBalance := currentPlayer.Wallet
	playerQuantity := currentPlayer.CoinQuantity

	totalSale := sellRequest.CurrentCoinPrice * sellRequest.Quantity
	if sellRequest.Quantity > playerQuantity {
		http.Error(w, "You do not have enough coin", http.StatusBadRequest)
		return
	} else {
		currentPlayer.Wallet = playerBalance + totalSale
		currentPlayer.CoinQuantity = playerQuantity - sellRequest.Quantity
		cacheStore.Set(currentPlayer.Name, currentPlayer, cache.DefaultExpiration)

		res := &CoinResponse{
			Player:        sellRequest.Player,
			CoinQuantity:  currentPlayer.CoinQuantity,
			WalletBalance: currentPlayer.Wallet,
		}
		channel.Publish(ctx, "sell", res)
		// writeResponse(res, w, r)
	}
}

func writeResponse(resValue interface{}, w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(resValue)
	if err != nil {
		ErrorLogger.Println(err)
		http.Error(w, "Json marshal error", http.StatusInternalServerError)
	}
	// Write response to http
	_, err = w.Write(res)
	if err != nil {
		ErrorLogger.Println(err)
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
	}
}

func AblyClient() *ably.Realtime {
	client, err := ably.NewRealtime(ably.WithKey(config.EnvVariable("API_KEY")))
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
