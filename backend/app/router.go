package app

import (
	"encoding/json"
	"net/http"
	"stocker/config"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

var store sessions.Store

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
	api := (*r)
	encryptionKey, err := determineEncryptionKey()
	if err != nil {
		ErrorLogger.Println(err)
	}
	store = sessions.NewCookieStore(
		[]byte(config.EnvVariable("SECRET")),
		encryptionKey,
	)

	api.HandleFunc("/start", panicRecover(CreatePlayer)).
		Methods(http.MethodPost)

	api.HandleFunc("/balance/{player}", panicRecover(GetWalletBalance)).
		Methods(http.MethodGet)

	api.HandleFunc("/buy", panicRecover(BuyBitcoin)).
		Methods(http.MethodPost)

	api.HandleFunc("/sell", panicRecover(SellBitcoin)).
		Methods(http.MethodPost)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{config.EnvVariable("FRONTEND")},
		AllowCredentials: true,
	}).Handler(api)
	return handler

}

// panicRecover - recover endpoint from panic
func panicRecover(restart func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				ErrorLogger.Println(err)
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

func determineEncryptionKey() ([]byte, error) {
	sek := config.EnvVariable("ENCRYPTION_KEY")
	lek := len(sek)
	switch {
	case lek >= 0 && lek < 16, lek > 16 && lek < 24, lek > 24 && lek < 32:
		return nil, errors.Errorf("SESSION_ENCRYPTION_KEY needs to be either 16, 24 or 32 characters long or longer, was: %d", lek)
	case lek == 16, lek == 24, lek == 32:
		return []byte(sek), nil
	case lek > 32:
		return []byte(sek[0:32]), nil
	default:
		return nil, errors.New("invalid SESSION_ENCRYPTION_KEY: " + sek)
	}

}
