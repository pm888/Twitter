package main

import (
	"Twitter_like_application/internal/database/Mongodb"
	"Twitter_like_application/internal/server"
)

func main() {
	server.Server()
	Mongodb.MongoDB()
}
