package main

import (
	"net/http"
	"time"
)

// handleJob memunculkan goroutine baru untuk tiap request tanpa batas.
func (s *Server) handleJob(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(10 * time.Second)
		// pekerjaan panjang tanpa worker pool / limit
	}()
	w.WriteHeader(http.StatusAccepted)
}
