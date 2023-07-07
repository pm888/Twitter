package admin

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	ServerHost string
	ServerPort string
	DbHost     string
)

func LoadEnvFile() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	DbPort = os.Getenv("DB_PG_PORT")
	DbUser = os.Getenv("DB_PG_USER")
	DbPassword = os.Getenv("DB_PG_PASSWORD")
	DbName = os.Getenv("DB_PG_NAME")
	ServerHost = os.Getenv("SERVER_HOST")
	ServerPort = os.Getenv("SERVER_PORT")
	DbHost = os.Getenv("DB_PG_HOST")
	return nil
}
