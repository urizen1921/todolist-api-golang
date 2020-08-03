package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

//User is a struct that represents a user
type User struct {
	FullName string `json:"fullName"`
	UserName string `json:"userName"`
	Email string `json:"email"`
}

//Post is a struct that represents a single post
type Post struct {
	Title string `json:"title"`
	Body string `json:"body"`
	Author User `json:"author"`
}

var posts []Post = []Post{}


func main() {
	router := mux.NewRouter()

	router.HandleFunc("/add", addItem).Methods("POST")

	http.ListenAndServe(":5000", router)

}


func addItem(w http.ResponseWriter, r *http.Request) {

	//get Item value from the JSON body
	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)

	posts = append(posts, newPost)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(posts)

}