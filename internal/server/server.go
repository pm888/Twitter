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
	r.HandleFunc("/getstatistic", Serviceuser.GetStatistics).Methods("GET")
	r.HandleFunc("/addusers", Serviceuser.CreateUser).Methods("POST")
	r.HandleFunc("/login", Serviceuser.LoginUsers).Methods("POST")
	http.Handle("/logout", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.LogoutUser)))
	http.Handle("/editprofile", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.EditProfile)))
	//r.HandleFunc("/resetpassword", Serviceuser.ResetPassword)
	//r.HandleFunc("/followuser", Serviceuser.FollowUser).Methods("POST")
	//r.HandleFunc("/unfollowuser", Serviceuser.UnfollowUser).Methods("POST")
	//r.HandleFunc("replacemyprofile", Serviceuser.EditMyProfile).Methods("POST")
	http.Handle("/myaccaunt", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.GetCurrentProfile)))
	r.HandleFunc("/getfollowers", Serviceuser.GetFollowers).Methods("GET")
	r.HandleFunc("/getfollowing", Serviceuser.GetFollowing).Methods("GET")
	r.HandleFunc("/searchusers", Serviceuser.SearchUsers).Methods("GET")
	r.HandleFunc("/searchtweet", Tweets.SearchTweets).Methods("GET")
	r.HandleFunc("/getfollowingtweet", Tweets.GetFollowingTweets).Methods("GET")
	http.Handle("/createtweet", Serviceuser.AuthHandler(http.HandlerFunc(Tweets.CreateTweet)))
	r.HandleFunc("/tweets_gertweet", Tweets.GetTweet).Methods("GET")
	r.HandleFunc("/tweets_gerpopulartweet", Tweets.GetPopularTweets).Methods("GET")
	r.HandleFunc("/tweets_rettweer", Tweets.Retweet).Methods("POST")
	r.HandleFunc("/tweets_liketweet", Tweets.LikeTweet).Methods("POST")
	r.HandleFunc("/tweets_unliketweet", Tweets.UnlikeTweet).Methods("POST")
	//r.HandleFunc("/tweets_updatetweet", Tweets.UpdateTweet).Methods("POST")
	//r.HandleFunc("/tweets_deletetweet", Tweets.DeleteTweet).Methods("POST")
	http.Handle("/myaccount", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.GetCurrentProfile)))
	r.HandleFunc("/v1/users/{id}/account", Serviceuser.GetUserProfile).Methods("GET")
	r.HandleFunc("/v1/users/password", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.ResetPassword)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	http.ListenAndServe("localhost:8080", r)
}
