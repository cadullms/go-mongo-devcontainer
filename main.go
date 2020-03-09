package main

// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

import (
	// "os"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

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
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err := viper.MergeInConfig(); err != nil {
		fmt.Println("Unable to load config from file.", err, "All config from env vars now.")
	}
	viper.SetConfigName("secrets")
	if err := viper.MergeInConfig(); err != nil {
		fmt.Println("Unable to load secrets from file.", err, "All secrets from env vars now.")
	}
	viper.AutomaticEnv()
}

func connectDb(mongoUrl string) *mongo.Client {
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

func getToken() string {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_secret", viper.GetString("CLIENT_SECRET"))
	form.Add("client_id", viper.GetString("CLIENT_ID"))
	form.Add("resource", viper.GetString("TARGET_RESOURCE_APP_ID"))
	body := strings.NewReader(form.Encode())
	response, _ := http.Post("https://login.microsoftonline.com/" + viper.GetString("TENANT_ID") + "/oauth2/token", "application/x-www-form-urlencoded", body)
	rawResult, _ := ioutil.ReadAll(response.Body)
	stringResult := string(rawResult)
	//var result map[string]interface{}
	var result map[string]string
	json.Unmarshal([]byte(stringResult), &result)
	return result["access_token"]
}

func getLessons() string {
	token := getToken()
	client := http.Client{}
	request, _ := http.NewRequest("GET", "https://ninaapp.carstenduellmann.de/api/lessons", nil)
	request.Header.Add("Authorization", "Bearer "+token)
	response, _ := client.Do(request)
	rawResult, _ := ioutil.ReadAll(response.Body)
	stringResult := string(rawResult)
	return stringResult
}

func main() {
	fmt.Println("Starting...")
	initConf()

	fmt.Println(getLessons())

	mongoUrl := viper.GetString("MONGO_URL")
	fmt.Println("Mongo url is", mongoUrl)
	client := connectDb(mongoUrl)
	collection := client.Database("test").Collection("attempts")

	attempt1 := Attempt{"user1", "lesson1"}
	attempt2 := Attempt{"user1", "lesson2"}

	collection.InsertOne(context.TODO(), attempt1)
	collection.InsertOne(context.TODO(), attempt2)
	fmt.Println("Inserted two docs")

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

}
