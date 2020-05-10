package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB
var err error

func main() {
	// Open database
	db, err = sql.Open("postgres", "postgresql://localhost:5432/ims")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	// Create endpoints
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

	// Start server
	http.ListenAndServe(":8000", router)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []Post

	// Get rows
	result, err := db.Query("SELECT id, title from posts")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	// Create sql statement
	stmt, err := db.Prepare("INSERT INTO posts(title) VALUES(?)") // ? - is a parameter placeholder for mysql
	if err != nil {
		panic(err.Error())
	}

	// Read from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	// Store request body as a map and get title
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]

	// Make query with title as a value
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New post was created")
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	// Get row
	result, err := db.Query("SELECT id, title from posts WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	var post Post

	for result.Next() {
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(post)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create sql statement
	stmt, err := db.Prepare("UPDATE posts SET title = ? WHERE id = ?") // ? - is a parameter placeholder for mysql
	if err != nil {
		panic(err.Error())
	}

	// Read from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	// Store request body as a map and get title
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]

	params := mux.Vars(r)

	_, err = stmt.Exec(title, params["id"])

	fmt.Fprintf(w, "Post %s was updated", params["id"])
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Delete row
	stmt, err := db.Prepare("DELETE from posts WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Post %s was deleted", params["id"])
}
