package server

import (
	Serviceuser "Twitter_like_application/cmd/internal/database"
	"Twitter_like_application/cmd/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

func Server() {
	r := mux.NewRouter()
	r.HandleFunc("/tweets", services.GetTweets).Methods("GET")
	r.HandleFunc("/tweets", services.CreateTweet).Methods("POST")
	r.HandleFunc("/addusers", Serviceuser.CreateUser).Methods("POST")
	http.ListenAndServe("localhost:8080", r)
}
