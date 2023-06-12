package server

import (
	Tweets "Twitter_like_application/internal/tweets"
	Serviceuser "Twitter_like_application/internal/users"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Server() {
	r := mux.NewRouter()
	fmt.Println("Server was run", "localhost:8080")
	http.ListenAndServe("localhost:8080", r)
	//r.HandleFunc("/getstatistic", Serviceuser.GetStatistics).Methods("GET")
	//r.HandleFunc("/resetpassword", Serviceuser.ResetPassword)
	//r.HandleFunc("/followuser", Serviceuser.FollowUser).Methods("POST")
	//r.HandleFunc("/unfollowuser", Serviceuser.UnfollowUser).Methods("POST")
	//r.HandleFunc("replacemyprofile", Serviceuser.EditMyProfile).Methods("POST")
	//r.HandleFunc("/getfollowers", Serviceuser.GetFollowers).Methods("GET")
	//r.HandleFunc("/getfollowing", Serviceuser.GetFollowing).Methods("GET")
	//r.HandleFunc("/searchusers", Serviceuser.SearchUsers).Methods("GET")
	//r.HandleFunc("/searchtweet", Tweets.SearchTweets).Methods("GET")
	//r.HandleFunc("/getfollowingtweet", Tweets.GetFollowingTweets).Methods("GET")
	//r.HandleFunc("/tweets_gertweet", Tweets.GetTweet).Methods("GET")
	//r.HandleFunc("/tweets_gerpopulartweet", Tweets.GetPopularTweets).Methods("GET")
	//r.HandleFunc("/tweets_rettweer", Tweets.Retweet).Methods("POST")
	//r.HandleFunc("/tweets_liketweet", Tweets.LikeTweet).Methods("POST")
	//r.HandleFunc("/tweets_unliketweet", Tweets.UnlikeTweet).Methods("POST")
	//r.HandleFunc("/tweets_updatetweet", Tweets.UpdateTweet).Methods("POST")
	//r.HandleFunc("/tweets_deletetweet", Tweets.DeleteTweet).Methods("POST")
	r.HandleFunc("/v1/users", Serviceuser.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/login", Serviceuser.LoginUsers).Methods(http.MethodPost)
	http.Handle("/v1/users/logout", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.LogoutUser)))
	http.Handle("/v1/users/{id}", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.GetCurrentProfile)))
	r.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.EditProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.CreateTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.EditTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
}
