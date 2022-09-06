package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/EleisonC/K-Dogo/configs"
	"github.com/EleisonC/K-Dogo/routes"
)

func main(){
	r := mux.NewRouter()
	routes.RegisterDogRoutes(r)
	http.Handle("/", r)
	// start the Database connection
	configs.ConnectDB()
	log.Fatal(http.ListenAndServe(":8082", r))
}