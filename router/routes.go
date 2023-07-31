package router

import (
	"github.com/durgesh730/authenticationInGo/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/players", controllers.CreatePlayer).Methods("POST")
	router.HandleFunc("/players/{id}", controllers.UpdateNameandScore).Methods("PUT")
	router.HandleFunc("/players/{id}", controllers.DeleteOnePlayer).Methods("DELETE")
	router.HandleFunc("/players/random", controllers.GetRandomPlayer).Methods("GET")
	router.HandleFunc("/players", controllers.GetAllPlayers).Methods("GET")
	router.HandleFunc("/players/rank/{val}", controllers.GetRankedPlayer).Methods("GET")

	return router
}