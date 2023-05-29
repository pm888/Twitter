package main

import (
	"Twitter_like_application/internal/database/Mongodb"
	"Twitter_like_application/internal/server"
	"fmt"
)

type ServiceMongoDb struct {
	DB interface{}
}

func main() {
	server.Server()
	var service = &ServiceMongoDb{
		DB: &Mongodb.MongoDB{},
	}
	err := service.ConnectPostgresql()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
