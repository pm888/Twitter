package server

import (
	Tweets "Twitter_like_application/internal/tweets"
	Serviceuser "Twitter_like_application/internal/users"
	"github.com/gorilla/mux"
	"net/http"
)

func Server() {
	r := mux.NewRouter()
	r.HandleFunc("/addusers", Serviceuser.CreateUser).Methods("POST")
	r.HandleFunc("/login", Serviceuser.LoginUsers).Methods("POST")
	r.HandleFunc("/logout", Serviceuser.LogoutUser)
	r.HandleFunc("replacemyprofile", Serviceuser.EditMyProfile).Methods("POST")
	r.HandleFunc("/myaccount", Serviceuser.GetUserProfile).Methods("GET")
	r.HandleFunc("/tweets", Tweets.CreateTweet).Methods("POST")
	r.HandleFunc("/tweets", Tweets.GetTweet).Methods("GET")
	r.HandleFunc("/tweets", Tweets.UpdateTweet).Methods("POST")
	r.HandleFunc("/tweets", Tweets.DeleteTweet).Methods("POST")
	//r.HandleFunc("/deleteuser", Serviceuser.DeleteUser).Methods("POST")
	//r.HandleFunc("/following", Serviceuser.Following).Methods("POST")
	//r.HandleFunc("/home", Serviceuser.Home)
	//r.HandleFunc("/resetpassword", Serviceuser.ResetPassword).Methods("POST")
	//r.HandleFunc("/liketweet", services.LikeTweet)
	http.ListenAndServe("localhost:8080", r)
}
