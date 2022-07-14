package main

import (
	"log"
	"net/http"

	"authentication.com/controllers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", controllers.HandleHome).Methods("GET")
	router.HandleFunc("/login", controllers.HandleLogin).Methods("POST")
	router.HandleFunc("/authenticated", controllers.HandleAuthentication).Methods("GET")
	router.HandleFunc("/refresh-token", controllers.HandleRefreshToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":1991", router))
}
