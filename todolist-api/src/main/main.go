package main

import (
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

//Task is a struct that represents a single task
type Task struct {
	gorm.Model
	Task string `json:"task"`
	Complete bool `json:"complete"`
}


func initialMigration() {


	db, err := gorm.Open("mysql", "urizen1921:1921urizen@tcp(127.0.0.1:3306)/tasktestgolang?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Task{})
}


func main() {

	initialMigration()


	router := mux.NewRouter()

	router.HandleFunc("/tasks", addTask).Methods("POST")

	router.HandleFunc("/tasks", getAllTasks).Methods("GET")

	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")

	router.HandleFunc("/tasks", updateTask).Methods("PUT")

	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	c := cors.New(cors.Options{
	AllowedOrigins: []string{"http://localhost:3000"},
	AllowCredentials: true,
	})

	handler := c.Handler(router)

	http.ListenAndServe(":5000", handler)

}

func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func getTask(w http.ResponseWriter, r *http.Request) {

	//get the id of the post from the route parameter
	var idParam string = mux.Vars(r)["id"]
	id, error := strconv.Atoi(idParam)

	if error != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer."))
		return
	}

	db, err := gorm.Open("mysql", "urizen1921:1921urizen@tcp(127.0.0.1:3306)/tasktestgolang?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	var task Task
	db.First(&task, id)

	if task.ID == 0 {
		w.WriteHeader(404)
		w.Write([]byte("No task found with specified ID"))
		return
	}

	fmt.Println("{}", task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

func updateTask(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open("mysql", "urizen1921:1921urizen@tcp(127.0.0.1:3306)/tasktestgolang?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	var updatedTask Task
	json.NewDecoder(r.Body).Decode(&updatedTask)
	db.Save(&updatedTask)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)

}

func getAllTasks(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open("mysql", "urizen1921:1921urizen@tcp(127.0.0.1:3306)/tasktestgolang?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	var tasks []Task
	db.Find(&tasks)
	fmt.Println("{}", tasks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}

func addTask(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open("mysql", "urizen1921:1921urizen@tcp(127.0.0.1:3306)/tasktestgolang?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect to the database")
	}

	defer db.Close()

	//get Item value from the JSON body
	var newTask Task
	json.NewDecoder(r.Body).Decode(&newTask)

	fmt.Println(newTask.Task);

	db.Create(&Task{
		Task: newTask.Task,
		IsComplete: newTask.IsComplete})
	fmt.Fprintf(w, "New Task Successfully Created")


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	/*
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	*/


	db, err := gorm.Open("mysql", "urizen1921:1921urizen@tcp(127.0.0.1:3306)/tasktestgolang?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	var idParam string = mux.Vars(r)["id"]
	fmt.Println(mux.Vars(r))
	id, er := strconv.Atoi(idParam)

	if er != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer."))
		return
	}

	var task Task
	db.First(&task, id)

	if task.ID == 0 {
		w.WriteHeader(404)
		w.Write([]byte("No task found with specified ID"))
		return
	}

	db.Delete(&task)

	fmt.Fprintf(w, "Successfully Deleted Task")

	w.WriteHeader(200)

}