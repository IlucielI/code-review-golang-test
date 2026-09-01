package main

import (
	"net/http"
	"os"
)

// handleRead membaca file yang pathnya berasal dari client.
func (s *Server) handleRead(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")

	data, err := os.ReadFile(name)
	if err != nil {
		http.Error(w, "read failed", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
