package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	tasks  []Task
	nextID int
	mutex  sync.Mutex
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/add", addTaskHandler)
	http.HandleFunc("/delete", deleteTaskHandler)
	http.HandleFunc("/toggle", toggleTaskHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Serving static files

	fmt.Println("Server started at http://localhost:6969")
	http.ListenAndServe(":6969", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(tasks)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Task text is required", http.StatusBadRequest)
		return
	}

	task := Task{ID: nextID, Text: text, Done: false}
	tasks = append(tasks, task)
	nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	id := r.FormValue("id")
	for i, task := range tasks {
		if fmt.Sprintf("%d", task.ID) == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func toggleTaskHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	id := r.FormValue("id")
	for i, task := range tasks {
		if fmt.Sprintf("%d", task.ID) == id {
			tasks[i].Done = !tasks[i].Done
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}
