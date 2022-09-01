package routes


import (
	"github.com/gorilla/mux"
	"github.com/EleisonC/K-Dogo/controllers"
)

var RegisterDogRoutes = func(router *mux.Router) {
	router.HandleFunc("{ownerId}/createdog", controllers.CreateDog).Methods("POST")
	router.HandleFunc("/getdog/{dogId}", controllers.GetDogById).Methods("GET")
	router.HandleFunc("/getdogs", controllers.GetAllDogs).Methods("GET")
	router.HandleFunc("/deletedog/{dogId}", controllers.DeleteDog).Methods("DELETE")
	router.HandleFunc("/updatedog/{dogId}", controllers.UpdateDog).Methods("PUT")
}