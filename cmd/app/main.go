package main

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/server"
	"fmt"
)

type ServiceMongoDb struct {
	DB interface{}
}

func main() {
	err := pg.ConnectPostgresql()
	fmt.Println(err)
	//if err := migrations.Run(pg.DB); err != nil {
	//	fmt.Println("running migrations", err)
	//}
	server.Server()

}
