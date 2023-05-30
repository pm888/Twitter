package main

import (
	"Twitter_like_application/internal/server"
)

type ServiceMongoDb struct {
	DB interface{}
}

func main() {
	server.Server()
}
