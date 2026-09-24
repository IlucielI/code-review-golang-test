package main

import (
	"io"
	"net/http"
)

// handleFetch mengambil URL arbitrary yang dikirim client.
func (s *Server) handleFetch(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "fetch failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}
