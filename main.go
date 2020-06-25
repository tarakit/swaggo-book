package main

// @title Book API
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/RestAPI/helper"
	"github.com/RestAPI/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

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

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// upload of 10MB files
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println(err)
		return
		// return "", err
	}
	defer file.Close()

	// generate uuid concatenated with chosen file name
	id, er := uuid.New()
	if er != nil {
		return
	}
	FILENAME := fmt.Sprint(id, handler.Filename)

	// f, err := os.OpenFile("C:/Users/kitta/go/src/github.com/RestAPI/photos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	f, err := os.OpenFile("./photos/"+FILENAME, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
		// return "", err
	}
	defer f.Close()
	io.Copy(f, file)

}

func main() {

	// init Router
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	r.HandleFunc("/api/upload", uploadFile).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", r))
}
