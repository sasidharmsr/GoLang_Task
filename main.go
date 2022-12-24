package main

import (
	"Task/db"
	"Task/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter() // Intiate a Router

	//run database
	db.ConnectDB()

	dbName := "go"

	// Routes
	router.HandleFunc("/user", handlers.PostUser(dbName)).Methods("POST")
	router.HandleFunc("/user", handlers.GetAllUsers(dbName)).Methods("GET")
	router.HandleFunc("/user/{id}", handlers.GetSingleUser(dbName)).Methods("GET")
	router.HandleFunc("/user/{id}", handlers.UpdateUser(dbName)).Methods("PUT")
	router.HandleFunc("/user/{id}", handlers.DeleteUser(dbName)).Methods("DELETE")

	router.HandleFunc("/user/follow/{id}", handlers.FollowUser(dbName)).Methods("POST")
	router.HandleFunc("/user/following/{id}", handlers.GetFollowingofUser(dbName)).Methods("GET")
	router.HandleFunc("/user/followers/{id}", handlers.GetFollowersofUser(dbName)).Methods("GET")

	router.HandleFunc("/user/getnearusers/{id}", handlers.GetNearByUsers(dbName)).Methods("GET")
	http.ListenAndServe(":8090", router)
}
