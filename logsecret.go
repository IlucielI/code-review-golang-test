package main

import (
	"log"
	"net/http"
)

// handleSignup mencatat password mentah ke log.
func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	log.Printf("new user %s with password %s", username, password)
	w.WriteHeader(http.StatusCreated)
}
