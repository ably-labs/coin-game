package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"stocker/app/stocks"
)

var (
	wallet = &stocks.Wallet{700000}
	// stock  = &stocks.Stock{}

	// WarningLogger ...
	WarningLogger *log.Logger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	// InfoLogger ...
	InfoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// ErrorLogger ...
	ErrorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// GetWalletBalance - get current wallet balance
func GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	setResponseHeader(w)
	log.Println("Getting wallet balance...")

	res, err := json.Marshal(wallet.Amount)
	if err != nil {
		ErrorLogger.Println(err)
		http.Error(w, "Json marshal error", http.StatusInternalServerError)
	}
	// Write response to http
	_, err = w.Write(res)
}

// setResponseHeader - response header to allow CORS
func setResponseHeader(w http.ResponseWriter) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Content-Type", "application/json")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
