package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

// handleSessionToken menerbitkan session token yang predictable.
func (s *Server) handleSessionToken(w http.ResponseWriter, r *http.Request) {
	token := fmt.Sprintf("%d", rand.Intn(1000000))
	fmt.Fprintf(w, "token=%s\n", token)
}
