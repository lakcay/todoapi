package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"math/rand"
	"github.com/gorilla/mux"
)

// ToDo struct (Model)
type ToDo struct {
	ID     string  `json:"id"`
	Title   string  `json:"title"`
	Description string `json:"description"`
	Tags     [ ]string `json:"tags"`
	DueDate   time.Time `json:"date"`
	Completed bool `json:"completed"`
	NotificationSent bool `json:"notificationsent"`
}

// Init todos var as a slice ToDo struct
var todos []ToDo

const DUE_DATE_CHECK_PERIOD = 10

// Get all todos
func getToDos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Get single todo
func getToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through todos and find one with the id from the params
	for _, item := range todos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&ToDo{})
}

// Get single todo by title
func getToDoByTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through todos and find one with the title from the params
	for _, item := range todos {
		if item.Title == params["title"] {
			json.NewEncoder(w).Encode(item)
			//return
		}
	}
	json.NewEncoder(w)
}


// Get single todo by description
func getToDoByDescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through todos and find one with the description from the params
	for _, item := range todos {
		if item.Description == params["description"] {
			json.NewEncoder(w).Encode(item)
			//return
		}
	}
	json.NewEncoder(w)
}

// Get  todo by tag
func getToDoByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through todos and find one with the tag from the params
	for _, item := range todos {
		for _, tag := range item.Tags {
			if tag == params["tag"] {
				json.NewEncoder(w).Encode(item)
				break;
			}
		}
	}
	json.NewEncoder(w)
}

// Get single todo by due date
func getToDoByDueDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through todos and find one with the date from the params
	for _, item := range todos {
		if item.DueDate.Format(time.RFC3339Nano) == params["date"] {
			json.NewEncoder(w).Encode(item)
		}
	}
	json.NewEncoder(w)
}

//Generating Id
func generateUUID() string{
b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
    log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
    b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// Add new todo
func createToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo ToDo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = generateUUID()
	todo.Completed = false
	todo.NotificationSent = false
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

// Update todo
func updateToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			var todo ToDo
			_ = json.NewDecoder(r.Body).Decode(&todo)
			todo.ID = params["id"]
			todos = append(todos, todo)
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
}

// Delete todo
func deleteToDo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(todos)
}

func deleteAllToDos(w http.ResponseWriter, r *http.Request) {
	todos = nil 
}

//this function is going to prepare a callback function according to client information (ios/android url etc)
func getNotificationCallbackFunction(itemUUID string, title string) func() int {
	return func() int {
	    fmt.Println("Item due date arrived. this is a push notification message for client who created item :" + title + " "+ itemUUID)	
        return 1
	}
}
//this function make a call to client for push notification
func pushNotification(notificationCallbackFunction func() int){
	notificationCallbackFunction()
}

func periodicDueDateChecker(tick time.Time){
	for index, item := range todos {
		if !item.NotificationSent  && (tick.Equal(item.DueDate) || tick.After(item.DueDate)) {
			// the due date arrives, we will send a push notification to the user
			//first prepare notification callback function according to client info
			notificationCallbackFunction := getNotificationCallbackFunction(item.ID, item.Title)
			//pass notificationCallbackFunction as a parameter to push function
			pushNotification(notificationCallbackFunction)
			todos[index].NotificationSent = true
		}
	}   
    
}

// Main function
func main() {
	
    ticker := time.NewTicker(DUE_DATE_CHECK_PERIOD * time.Second)
    go func(){
        for t := range ticker.C {
            //Call the periodic function here.
            periodicDueDateChecker(t)
        }
    }()
	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/todos", getToDos).Methods("GET")
	r.HandleFunc("/todos/", getToDos).Methods("GET")
	r.HandleFunc("/todos/{id}", getToDo).Methods("GET")
	r.HandleFunc("/todos/title/{title}", getToDoByTitle).Methods("GET")
	r.HandleFunc("/todos/description/{description}", getToDoByDescription).Methods("GET")
	r.HandleFunc("/todos/duedate/{date}", getToDoByDueDate).Methods("GET")
	r.HandleFunc("/todos/tag/{tag}", getToDoByTag).Methods("GET")
	r.HandleFunc("/todos",  createToDo).Methods("POST")
	r.HandleFunc("/todos/{id}", updateToDo).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteToDo).Methods("DELETE")
	r.HandleFunc("/todos", deleteAllToDos).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))	
}


