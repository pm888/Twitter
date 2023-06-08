package main

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/server"
)

type ServiceMongoDb struct {
	DB interface{}
}

func main() {
	pg.ConnectPostgresql()
	server.Server()
}
