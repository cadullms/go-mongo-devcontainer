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
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Unable to load config from file.", err, "All config from env vars now.")
	}
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
	tenantId := "49c22610-756e-4632-a268-0712e0fc1ef5"
	clientId := "b97356b1-23db-4ecc-b977-81b25ee0657f"
	clientSecret := ""
	targetResourceAppId := "42bdba43-a7fe-4c9c-a2d4-44857edee58e"
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_secret", clientSecret)
	form.Add("client_id", clientId)
	form.Add("resource", targetResourceAppId)
	body := strings.NewReader(form.Encode())
	response, _ := http.Post("https://login.microsoftonline.com/"+tenantId+"/oauth2/token", "application/x-www-form-urlencoded", body)
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
