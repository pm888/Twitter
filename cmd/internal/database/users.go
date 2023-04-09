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
	Following []string
	Followers []int
}

type Tweet struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
