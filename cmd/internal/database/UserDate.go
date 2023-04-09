package Serviceuser

import "fmt"

type UserDataSTR struct {
	UserData map[int]*Users
}

var (
	UserDate = make(map[int]*Users)
	counter  = 0
)

func Put(u *Users) {
	for _, user := range UserDate {
		if user.Email == u.Email {
			fmt.Println("This user alredy was added ")

		}
	}
	counter++
	u.ID = counter
	UserDate[u.ID] = u
}
