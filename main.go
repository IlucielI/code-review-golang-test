package main

import (
	"database/sql"
	"log"
	"net/http"
)

type User struct {
	ID   int
	Name string
}

type Server struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/app")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := &Server{db: db}

	http.HandleFunc("/posts/search", s.handleSearch)
	http.HandleFunc("/posts", s.handleList)
	http.HandleFunc("/posts/render", s.handleRender)
	http.HandleFunc("/reports/export", s.handleExport)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
