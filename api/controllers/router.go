package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// HandleRequests creates the router for the applicaation 
func HandleRequests(port string) {

	fmt.Println("Running on port:" + port)

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/authorised", authorised).Methods("GET")
	myRouter.HandleFunc("/signup", signUp).Methods("POST")

	cors := handlers.AllowedOrigins([]string{"*"})
	
	log.Fatal(http.ListenAndServe(":" + port, handlers.CORS(cors)(myRouter)))
}