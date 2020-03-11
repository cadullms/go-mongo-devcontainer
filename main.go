package main

// https://medium.com/@bnprashanth256/reading-configuration-files-and-environment-variables-in-go-golang-c2607f912b63
// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"dev.azure.com/go-mongo/model"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spf13/viper"
)

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

func crashIt(err error) {
	log.Println(err)
	os.Exit(1) // Note: Some debuggers still claim the process exited with code 0 after this. Try with go run or with the compiled version to verify non zero exit code!
}

func connectDb() (*mongo.Client, error) {
	mongoUrl := viper.GetString("MONGO_URL")
	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func getToken() (string, error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_secret", viper.GetString("CLIENT_SECRET"))
	form.Add("client_id", viper.GetString("CLIENT_ID"))
	form.Add("resource", viper.GetString("TARGET_RESOURCE_APP_ID"))
	body := strings.NewReader(form.Encode())
	response, httperr := http.Post("https://login.microsoftonline.com/"+viper.GetString("TENANT_ID")+"/oauth2/token", "application/x-www-form-urlencoded", body)
	if httperr != nil {
		return "", httperr
	}

	rawResult, ioerr := ioutil.ReadAll(response.Body)
	if ioerr != nil {
		return "", ioerr
	}

	stringResult := string(rawResult)
	//var result map[string]interface{}
	var result map[string]string
	jsonerror := json.Unmarshal([]byte(stringResult), &result)
	if jsonerror != nil {
		return "", jsonerror
	}

	return result["access_token"], nil
}

func get(url string) ([]byte, error) {

	var err error
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, errors.New("Call to url " + url + " resulted in " + response.Status + ".")
	}

	rawResult, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return rawResult, nil
}

func getLessons() ([]model.Lesson, error) {
	url := viper.GetString("API_BASE_URL") + "api/lessons"
	rawResult, err := get(url)
	if err != nil {
		return nil, err
	}

	var lessons []model.Lesson
	err = json.Unmarshal(rawResult, &lessons)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func getAttempts() ([]model.Attempt, error) {
	url := viper.GetString("API_BASE_URL") + "session/attempts/all"
	rawResult, err := get(url)
	if err != nil {
		return nil, err
	}

	var attempts []model.Attempt
	err = json.Unmarshal(rawResult, &attempts)
	if err != nil {
		return nil, err
	}

	return attempts, nil
}

func main() {
	fmt.Println("Starting...")
	initConf()

	lessons, err := getLessons()
	if err != nil {
		crashIt(err)
	}
	fmt.Println("Got ", len(lessons), " lessons.")

	attempts, err := getAttempts()
	if err != nil {
		crashIt(err)
	}
	fmt.Println("Got ", len(attempts), " attempts.")

    // TODO: get all user ids and/or filter attempts on this user

	for _, lesson := range lessons {
		model.ScoreAttemptsForLesson(&lesson, "live.com#cadull@hotmail.de", &attempts)
	}

	client, err := connectDb()
	if err != nil {
		crashIt(err)
	}

	db := client.Database("test")
	collection := db.Collection("lessons")

	_, err = collection.InsertOne(context.TODO(), lessons[0])
	if err != nil {
		crashIt(err)
	}

	fmt.Println("Inserted.")

	err = client.Disconnect(context.TODO())
	if err != nil {
		crashIt(err)
	}
}
