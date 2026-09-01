package main

import "net/http"

var sharedStats = make(map[string]int)

// handleWarm fires a background writer mutating a shared map.
func (s *Server) handleWarm(w http.ResponseWriter, r *http.Request) {
	go func() {
		for i := 0; i < 100; i++ {
			sharedStats["counter"] = i
		}
	}()
	w.WriteHeader(http.StatusAccepted)
}
