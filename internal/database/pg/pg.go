package pg

import (
	"database/sql"
	"fmt"
)

const (
	postgresUsername = "postgres"
	postgresPassword = "postgrespw"
	postgresDBip     = "localhost"
	postgresDBName   = "tweeter"
	portPG           = "55000"
)

var DB *sql.DB

type ServicePostgresql struct {
	DB *sql.DB
}

func ConnectPostgresql() error {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:55000/%s?sslmode=disable", postgresUsername, postgresPassword, postgresDBip, postgresDBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("error connect BD")
		return err
	}

	DB = db
	fmt.Println("**** PG ran.... ****")

	return nil
}
