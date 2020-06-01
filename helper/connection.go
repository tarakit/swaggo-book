package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

func ConnectDB() *mongo.Collection {

	// Set Client options
	clientOption := options.Client().ApplyURI("mongodb+srv://postgres:postgres@cluster0-rcysg.mongodb.net/test?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!!")

	collection := client.Database("my_new_db").Collection("Book")

	return collection
}

// Error Response
type ErrorResponse struct {
	StatusCode  int    `json:"code"`
	CodeMessage string `json:"message"`
}

// this function help to get error model and response to the webpage
func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		CodeMessage: err.Error(),
		StatusCode:  http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	// check error message
	w.Write(message)
}
