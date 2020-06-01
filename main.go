package main

import (
	"RestAPI/helper"
	"RestAPI/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

)

// get all data from Book Collection
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create a context
	ctx := context.TODO()

	var books []models.Book

	// create and retrieve connection from helper
	collection := helper.ConnectDB()

	cur, err := collection.Find(ctx, bson.M{})

	// check weither the cursor of indexing on the collection will got error reponse back
	if err != nil {
		helper.GetError(err, w)
		return
	}

	// defer : is used to close the conneciton when it is done
	defer cur.Close(ctx)

	// travel through the collection of data
	for cur.Next(ctx) {

		var book models.Book

		err := cur.Decode(&book) // decode to deserialized the books from Collection of MongoDB
		if err != nil {
			log.Fatal(err)
		}

		books = append(books, book)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// encode to serialize object of book
	json.NewEncoder(w).Encode(books)
}

func main() {

	// init Router
	r := mux.NewRouter()

	r.HandleFunc("/api/books", getBooks).Methods("GET")

	log.Fatal(http.ListenAndServe(":8002", r))
}
