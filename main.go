package main

// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

import (
	// "os"
	"context"
	"fmt"
	"log"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spf13/viper"
)

type Attempt struct {
	UserId   string
	LessonId string
}

type Lesson struct {
	LessonId string
}

func initConf() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Unable to load config from file.", err, "All config from env vars now.")
	}
}

func connectDb(mongoUrl string) (*mongo.Client) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoUrl)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err) // ends the program
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func main() {
	fmt.Println("Starting...")
	initConf()
	mongoUrl := viper.GetString("MONGO_URL")
	fmt.Println("Mongo url is", mongoUrl)
	client := connectDb(mongoUrl)
	collection := client.Database("test").Collection("attempts")

	attempt1 := Attempt{"user1","lesson1"}
	attempt2 := Attempt{"user1","lesson2"}

	collection.InsertOne(context.TODO(), attempt1)
	collection.InsertOne(context.TODO(), attempt2)
	fmt.Println("Inserted two docs")
	
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

}
