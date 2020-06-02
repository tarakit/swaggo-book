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
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book models.Book

	collection := helper.ConnectDB()

	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&book)

	err := collection.FindOneAndDelete(context.TODO(), filter).Decode(&book)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r) // get the params with Mux

	id, _ := primitive.ObjectIDFromHex(params["id"])

	var book models.Book

	collection := helper.ConnectDB()

	filter := bson.M{"_id": id}

	// read update model from body requested
	_ = json.NewDecoder(r.Body).Decode(&book)

	update := bson.D{
		{"$set", bson.D{
			{"title", book.Title},
			{"author", book.Author}}, // put the Author object rather than the document of Author
		},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&book)

	if err != nil {
		helper.GetError(err, w)
		return
	}
	book.ID = id

	json.NewEncoder(w).Encode(book)
}

// create book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	ctx := context.TODO()

	_ = json.NewDecoder(r.Body).Decode(&book)

	collection := helper.ConnectDB()

	result, err := collection.InsertOne(ctx, book)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book

	var params = mux.Vars(r) // get the params with Mux

	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := helper.ConnectDB()

	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&book)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// get all data from Book Collection
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create a context
	ctx := context.TODO()

	var books []models.Book

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
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
