package Mongodb

import (
	Serviceuser "Twitter_like_application/internal/database"
	"Twitter_like_application/internal/services"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

const (
	postgresUsername = "your_username"
	postgresPassword = "your_password"
	postgresDBName   = "your_dbname"
)

type ServiceMongoDb struct {
	DBClient *mongo.Client
	DB       *sql.DB
}

func (s *ServiceMongoDb) ConnectPostgresql() error {
	var w http.ResponseWriter

	connStr := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", postgresUsername, postgresPassword, postgresDBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	return nil
}
func (s ServiceMongoDb) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser Serviceuser.Users
	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Nickname == "" {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	query := `INSERT INTO users (name, password, email, nickname) VALUES ($1, $2, $3, $4) RETURNING id`
	err = s.DB.QueryRow(query, newUser.Name, newUser.Password, newUser.Email, newUser.Nickname).Scan(&newUser.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, "This user is already added", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userToken := services.CheckEmail(&newUser)
	newUser.EmailToken = userToken

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (s *ServiceMongoDb) LoginUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		usermail := r.FormValue("usermail")
		password := r.FormValue("password")

		query := "SELECT COUNT(*) FROM users WHERE email = $1 AND password = $2"
		var count int
		err := s.DB.QueryRow(query, usermail, password).Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			cookie := &http.Cookie{
				Name:  "session",
				Value: "authenticated",
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, nil)
		}
	}
}
func (s *ServiceMongoDb) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var deleteUser *Serviceuser.Users
	err := json.NewDecoder(r.Body).Decode(&deleteUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "your_connection_string_here")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = $1", deleteUser.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == "undefined_table" {
			http.Error(w, "Table not found", http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ServiceMongoDb) EditMyProfile(w http.ResponseWriter, r *http.Request) {
	var newuser *Serviceuser.ReplaceMyData

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newuser.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Exec(`
		UPDATE UserData SET 
			Name = $1,
			Email = $2,
			Nickname = $3,
			BirthDate = $4,
			Bio = $5,
			Password = $6,
			Location = $7
		WHERE id = $8`,
		newuser.NewName,
		newuser.NewEmail,
		newuser.NewNickname,
		newuser.NewBirthDate,
		newuser.NewBio,
		string(hashedNewPassword),
		newuser.NewLocation,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var userResPass Serviceuser.Users
	err := json.NewDecoder(r.Body).Decode(&userResPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	services.ResetPasswordPlusEmail(&userResPass)
}
func (s ServiceMongoDb) Following(w http.ResponseWriter, r *http.Request) {
	var user Serviceuser.FollowingForUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writer := user.Writer
	subscriber := user.Subscriber

	collection := s.DBClient.Database("your_database").Collection("your_collection")

	filterWriter := bson.M{"_id": writer}
	updateWriter := bson.M{"$addToSet": bson.M{"following": subscriber}}
	_, err = collection.UpdateOne(context.TODO(), filterWriter, updateWriter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filterSubscriber := bson.M{"_id": subscriber}
	updateSubscriber := bson.M{"$addToSet": bson.M{"followers": writer}}
	_, err = collection.UpdateOne(context.TODO(), filterSubscriber, updateSubscriber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
