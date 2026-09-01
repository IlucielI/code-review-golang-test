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
	http.HandleFunc("/admin/exec", s.handleExec)
	http.HandleFunc("/proxy/fetch", s.handleFetch)
	http.HandleFunc("/files/read", s.handleRead)
	http.HandleFunc("/files/save", s.handleSave)
	http.HandleFunc("/session/token", s.handleSessionToken)
	http.HandleFunc("/users/signup", s.handleSignup)
	http.HandleFunc("/echo", s.handleEcho)
	http.HandleFunc("/jobs/run", s.handleJob)
	http.HandleFunc("/health/blocks", s.handleBlock)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
