package main

import (
	"Twitter_like_application/internal/admin"
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/server"
	"Twitter_like_application/migrations"
	"fmt"
	"sync"
)

type ServiceMongoDb struct {
	DB interface{}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := admin.LoadEnvFile()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("**** env files is completed ****")
		}
		wg.Done()
	}()
	wg.Wait()
	err := pg.ConnectPostgresql()
	if err != nil {
		fmt.Println(err)
	}
	err = migrations.Run(pg.DB)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("**** running migrations ****", err)
	}
	err = server.Server()
	if err != nil {
		fmt.Println(err)
	}

}
