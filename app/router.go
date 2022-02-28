package app

import (
	"encoding/json"
	"net/http"
	handler "stocker/app/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	routePrefix = "/play"
)

// Router - router struct
type Router struct {
	*mux.Router
}

// NewRouter - new router instance
func NewRouter() *Router {
	return &Router{mux.NewRouter()}
}

// InitializeRoutes ...
func (r *Router) InitializeRoutes() http.Handler {
	api := (*r).PathPrefix(routePrefix).Subrouter()

	api.HandleFunc("/", panicRecover(handler.GetWalletBalance)).
		Methods(http.MethodGet)

	api.HandleFunc("/buy", panicRecover(handler.BuyBitcoin)).
		Methods(http.MethodGet)

	handler := cors.Default().Handler(api)
	return handler

}

// panicRecover - recover endpoint from panic
func panicRecover(restart func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				handler.ErrorLogger.Println(err)
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		restart(w, r)
	}
}
