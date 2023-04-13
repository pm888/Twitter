package Serviceuser

import (
	"time"
)

type Users struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Tweet
	Following []int
	Followers []int
}

type Tweet struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteUserST struct {
	UserIdDeleting int `json:"delete_id"`
}

type FollowingForUser struct {
	Writer     int `json:"writer"`
	Subscriber int `json:"subscriber"`
	Users
}
