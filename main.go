package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/router"
)

func main() {

	// connect database
	err := database.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to mongoDB", err)
		return
	}

	// connect router
	r := router.Router()

	//start the server
	fmt.Println("server linstening on port http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}
