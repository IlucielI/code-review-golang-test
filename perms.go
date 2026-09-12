package main

import (
	"net/http"
	"os"
)

// handleSave menulis data dengan permission world-writable.
func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	data := []byte(r.URL.Query().Get("data"))

	if err := os.WriteFile("/tmp/upload.txt", data, 0777); err != nil {
		http.Error(w, "write failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
