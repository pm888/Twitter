package Mongodb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type ConfigMongoDb struct {
	MongoUser     string `mapstructure:"mongoUser"`
	MongoPassword string `mapstructure:"mongoPassword"`
}

func ConnectPostgresql(w http.ResponseWriter) (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgresql://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	defer db.Close()
	return db, err
}

func LoadConfig(path string) (c ConfigMongoDb, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}

func readYaml() {
	config, err := LoadConfig(".")

	if err != nil {
		panic(fmt.Errorf("fatal error with config.yaml: %w", err))
	}
}

func MongoDB() {
	f := "name"
	config, err := LoadConfig(".")

	if err != nil {
		panic(fmt.Errorf("fatal error with config.yaml: %w", err))
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://:<>@cluster0.hleclhd.mongodb.net/?retryWrites=true&w=majority", f).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
