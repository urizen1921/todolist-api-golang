package main

import (
	"strconv"
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

	router.HandleFunc("/posts", addItem).Methods("POST")

	router.HandleFunc("/posts", getAllPosts).Methods("GET")

	router.HandleFunc("/posts/{id}", getPost).Methods("GET")

	router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")

	router.HandleFunc("/posts/{id}", patchPost).Methods("PATCH")

	http.ListenAndServe(":5000", router)

}

func getPost(w http.ResponseWriter, r *http.Request) {

	//get the id of the post from the route parameter
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		// there is an error
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer."))
		return
	}

	//error checking

	if id >= len(posts) {
		w.WriteHeader(404)
		w.Write([]byte("No post found with specified ID"))
		return
	}

	post := posts[id]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}

func updatePost(w http.ResponseWriter, r *http.Request) {

	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	if id >= len(posts) {
		w.WriteHeader(404)
		w.Write([]byte("No post found with specified ID"))
		return
	}

	//get the value from the JSON body

	var updatedPost Post
	json.NewDecoder(r.Body).Decode(&updatedPost)

	posts[id] = updatedPost

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPost)
}

func patchPost(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	if id >= len(posts) {
		w.WriteHeader(404)
		w.Write([]byte("No post found with specified ID"))
		return
	}

	post := &posts[id]
	json.NewDecoder(r.Body).Decode(post)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}


func addItem(w http.ResponseWriter, r *http.Request) {

	//get Item value from the JSON body
	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)

	posts = append(posts, newPost)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}