package routes

import (
	"github.com/gorilla/mux"
	"github.com/undersleep7x/cryo-project/handlers"
)

func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/ping", handlers.Ping).Methods("GET")         // ping route
	router.HandleFunc("/price", handlers.FetchPrices).Methods("GET") // route for sourcing pricing data from CoinGecko API
}
