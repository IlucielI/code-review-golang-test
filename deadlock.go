package main

import (
	"net/http"
)

// handleBlock mengirim ke unbuffered channel tanpa ada receiver.
func (s *Server) handleBlock(w http.ResponseWriter, r *http.Request) {
	ch := make(chan int)

	ch <- 1 // blocking selamanya; tidak ada penerima
	w.WriteHeader(http.StatusOK)
}
