package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"stocker/app/stocks"
	"strconv"
)

var (
	wallet = &stocks.Wallet{Amount: 700000}
	coin   = &stocks.Coin{}

	// WarningLogger ...
	WarningLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	// InfoLogger ...
	InfoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// ErrorLogger ...
	ErrorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// GetWalletBalance - get current wallet balance
func GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting wallet balance...")

	writeResponse(wallet.Amount, w, r)
}

func BuyBitcoin(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	price, err := strconv.ParseFloat(urlParams.Get("price"), 32)
	prices := float32(price)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	quantity, err := strconv.ParseFloat(urlParams.Get("qty"), 32)
	qty := float32(quantity)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	totalCost := prices * qty
	if totalCost > wallet.Amount {
		http.Error(w, "Insufficient wallet balance", http.StatusBadRequest)
		return
	} else {
		wallet.Subtract(totalCost)
		coin.AddQuantity(qty)

		res := &stocks.CoinResponse{
			Quantity:      coin.Quantity,
			WalletBalance: wallet.Amount,
		}
		writeResponse(res, w, r)

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
