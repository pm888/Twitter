package Serviceuser

type UserDataSTR struct {
	UserData map[int]*Users
}

var (
	UserDate = make(map[int]*Users)
)

func Put(u *Users) bool {
	for _, user := range UserDate {
		if user.Email == u.Email {
			return false

		}

	}
	u.ID = len(UserDate) + 1
	UserDate[u.ID] = u

	return true

}
