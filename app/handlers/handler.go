package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"stocker/app/stocks"
	"stocker/config"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	// "github.com/ably/ably-go/ably"
)

var (
	key     = []byte(config.EnvVariable("SECRET"))
	store   = sessions.NewCookieStore(key)
	channel = AblyClient().Channels.Get("stock")
	ctx     = context.Background()

	// WarningLogger ...
	WarningLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	// InfoLogger ...
	InfoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// ErrorLogger ...
	ErrorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("creating coin....")
	player := stocks.Player{}
	err := json.NewDecoder(r.Body).Decode(&player)
	player.Wallet = 700000
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, _ := store.Get(r, player.Name)
	session.Values["balance"] = player.Wallet
	session.Values["quantity"] = player.CoinQuantity
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(player.Name, "+n\name")
	// err = channel.Presence.Enter(ctx, player.Name)
	// if err != nil {
	// 	panic(err)
	// }

	// RegisterPresense(ctx, channel, player.Name)
	// clients, err := channel.Presence.Get(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(clients, "\n\n+++++")

	// for _, client := range clients {
	// 	fmt.Println("Present client:", client)
	// }
	// fmt.Println(clients, "======\nour clients")
	w.WriteHeader(http.StatusCreated)
}

// GetWalletBalance - get current wallet balance
func GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting wallet balance...")
	var res *stocks.CoinResponse
	player := mux.Vars(r)["player"]
	fmt.Println(player, "player name")
	session, err := store.Get(r, player)
	if err != nil {
		return
	}
	res = &stocks.CoinResponse{
		Player:        player,
		CoinQuantity:  session.Values["quantity"].(float32),
		WalletBalance: session.Values["balance"].(float32),
	}
	fmt.Println(res, "getting walltet")
	writeResponse(res, w, r)
}

func BuyBitcoin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("buying coin", r.Body)

	buyRequest := stocks.BuySellRequest{}
	err := json.NewDecoder(r.Body).Decode(&buyRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("player name", buyRequest.Player)
	player, err := store.Get(r, buyRequest.Player)
	if err != nil {
		fmt.Println(err, "error p")
		return
	}

	playerBalance := player.Values["balance"].(float32)
	playerQuantity := player.Values["quantity"].(float32)

	totalCost := buyRequest.CurrentCoinPrice * buyRequest.Quantity
	if totalCost > playerBalance {
		http.Error(w, "Insufficient wallet balance", http.StatusBadRequest)
		return
	} else {
		player.Values["balance"] = playerBalance - totalCost
		player.Values["quantity"] = playerQuantity + buyRequest.Quantity
		player.Save(r, w)

		res := &stocks.CoinResponse{
			Player:        buyRequest.Player,
			CoinQuantity:  player.Values["quantity"].(float32),
			WalletBalance: player.Values["balance"].(float32),
		}

		channel.Publish(ctx, "buy", res)

		if err != nil {
			panic(err)
		}
		// writeResponse(res, w, r)
	}

}

func SellBitcoin(w http.ResponseWriter, r *http.Request) {
	sellRequest := stocks.BuySellRequest{}
	err := json.NewDecoder(r.Body).Decode(&sellRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player, _ := store.Get(r, sellRequest.Player)
	playerBalance := player.Values["balance"].(float32)
	playerQuantity := player.Values["quantity"].(float32)

	totalSale := sellRequest.CurrentCoinPrice * sellRequest.Quantity
	if sellRequest.Quantity > playerQuantity {
		http.Error(w, "You do not have enough coin", http.StatusBadRequest)
		return
	} else {
		player.Values["balance"] = playerBalance + totalSale
		player.Values["quantity"] = playerQuantity - sellRequest.Quantity
		player.Save(r, w)

		res := &stocks.CoinResponse{
			Player:        sellRequest.Player,
			CoinQuantity:  player.Values["quantity"].(float32),
			WalletBalance: player.Values["balance"].(float32),
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
