package server

import (
	"Twitter_like_application/internal/database"
	"Twitter_like_application/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

func Server() {
	r := mux.NewRouter()
	//r.HandleFunc("/tweets", services.GeTweetsLast100).Methods("GET")
	r.HandleFunc("/tweets", services.CreateTweet).Methods("POST")
	r.HandleFunc("/addusers", Serviceuser.CreateUser).Methods("POST")
	r.HandleFunc("/deleteuser", Serviceuser.DeleteUser).Methods("POST")
	r.HandleFunc("/login", Serviceuser.LoginUsers).Methods("POST")
	r.HandleFunc("/following", Serviceuser.Following).Methods("POST")
	r.HandleFunc("replacemyprofile", Serviceuser.EditmyProfile).Methods("POST")
	r.HandleFunc("/myaccount", Serviceuser.ExploreMyaccaunt).Methods("GET")
	r.HandleFunc("/logout", Serviceuser.LogoutUser)
	r.HandleFunc("/home", Serviceuser.Home)
	r.HandleFunc("/resetpassword", Serviceuser.ResetPassword).Methods("POST")
	r.HandleFunc("/liketweet", services.LikeTweet)
	//r.HandleFunc("/getuser", Serviceuser.GetUser).Methods("GET")
	http.ListenAndServe("localhost:8080", r)
}
