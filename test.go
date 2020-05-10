package Course

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)



//Book Struct (Model)
type Book struct {
	ID    string `json:"ID"`
	Title string `json:"Title"`
}

//Init database and error vars
var db *sql.DB
var err error
var books []Book
//Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	result, err := db.Query("SELECT * from Books")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {

		var book Book

		err := result.Scan(&book.ID, &book.Title)
		if err != nil {
			panic(err.Error())
		}

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

//Get Single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //Get params

	result, err := db.Query("SELECT * FROM Books WHERE ID = @ID", sql.Named("ID", params["ID"]))
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var book Book

	for result.Next() {
		err := result.Scan(&book.ID, &book.Title)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(book)
}

//Create a New Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	ID := keyVal["ID"]
	newTitle := keyVal["Title"]

	result, err := db.Query("INSERT INTO Books(ID,Title) VALUES(@ID,@newTitle)", sql.Named("ID", ID), sql.Named("newTitle", newTitle))
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	fmt.Fprintf(w, "New book was created")

}

//Update Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["Title"]

	result, err := db.Query("UPDATE Books SET Title = @newTitle WHERE ID = @ID", sql.Named("newTitle", newTitle), sql.Named("ID", params["ID"]))
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	fmt.Fprintf(w, "Book with ID = %s was updated", params["ID"])

}

//Delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("DELETE FROM Books WHERE ID = @ID", sql.Named("ID", params["ID"]))
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	fmt.Fprintf(w, "Book with ID = %s was deleted", params["ID"])

}

func main() {
	db, err = sql.Open("postgres" , "postgresql://localhost:5432/ims")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	//Init Router
	r := mux.NewRouter()

	//Route Handlers/Endpoints
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{ID}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{ID}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{ID}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8001", r))
}
