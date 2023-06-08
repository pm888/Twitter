package pg

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

var DB *sql.DB

type ServicePostgresql struct {
	DBClient *mongo.Client
	DB       *sql.DB
}

func ConnectPostgresql() error {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", postgresUsername, postgresPassword, postgresDBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	DB = db

	return nil
}
