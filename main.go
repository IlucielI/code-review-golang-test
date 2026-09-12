package main

import (
	"database/sql"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/app")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskRepo := NewMemoryTaskRepository()
	taskHandler := NewTaskHandler(taskRepo)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/tasks", taskHandler.HandleTasks)
	http.HandleFunc("/tasks/", taskHandler.HandleTaskByID)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
