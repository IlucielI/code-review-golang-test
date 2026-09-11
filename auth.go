package main

import "net/http"

const adminToken = "s3cr3t-admin-token-2024"

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Admin-Token") != adminToken {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}
