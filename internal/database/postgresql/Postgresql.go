package Postgresql

import (
	"database/sql"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	postgresUsername = "your_username"
	postgresPassword = "your_password"
	postgresDBName   = "your_dbname"
)

type ServicePostgresql struct {
	DBClient *mongo.Client
	DB       *sql.DB
}

func (s *ServicePostgresql) ConnectPostgresql() error {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", postgresUsername, postgresPassword, postgresDBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	s.DB = db

	return nil
}

//func (s *ServicePostgresql) DeleteUser(w http.ResponseWriter, r *http.Request) {
//	var deleteUser *Serviceuser.Users
//	err := json.NewDecoder(r.Body).Decode(&deleteUser)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	db, err := sql.Open("postgres", "your_connection_string_here")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer db.Close()
//
//	_, err = db.Exec("DELETE FROM users WHERE id = $1", deleteUser.ID)
//	if err != nil {
//		pqErr, ok := err.(*pq.Error)
//		if ok && pqErr.Code.Name() == "undefined_table" {
//			http.Error(w, "Table not found", http.StatusInternalServerError)
//		} else {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

//func (s ServicePostgresql) Following(w http.ResponseWriter, r *http.Request) {
//	var user Serviceuser.FollowingForUser
//	err := json.NewDecoder(r.Body).Decode(&user)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	writer := user.Writer
//	subscriber := user.Subscriber
//
//	collection := s.DBClient.Database("your_database").Collection("your_collection")
//
//	filterWriter := bson.M{"_id": writer}
//	updateWriter := bson.M{"$addToSet": bson.M{"following": subscriber}}
//	_, err = collection.UpdateOne(context.TODO(), filterWriter, updateWriter)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	filterSubscriber := bson.M{"_id": subscriber}
//	updateSubscriber := bson.M{"$addToSet": bson.M{"followers": writer}}
//	_, err = collection.UpdateOne(context.TODO(), filterSubscriber, updateSubscriber)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}
