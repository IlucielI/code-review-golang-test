package main

import (
	"net/http"
)

// handleEcho merefleksikan input user ke response header.
func (s *Server) handleEcho(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	w.Header().Set("X-User", name)
	w.WriteHeader(http.StatusOK)
}
