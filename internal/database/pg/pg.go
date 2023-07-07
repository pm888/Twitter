package pg

import (
	"Twitter_like_application/internal/admin"
	"database/sql"
	"fmt"
)

var DB *sql.DB

type ServicePostgresql struct {
	DB *sql.DB
}

func ConnectPostgresql() error {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", admin.DbUser, admin.DbPassword, admin.DbHost, admin.DbPort, admin.DbName)
	db, err := sql.Open(fmt.Sprintf("%s", admin.DbUser), connStr)
	if err != nil {
		fmt.Println("error connect BD")
		return err
	}

	DB = db
	fmt.Println("**** PG ran.... ****")

	return nil
}
