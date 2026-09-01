package main

import (
	"net/http"
	"os/exec"
)

// handleExec menjalankan perintah shell yang dibangun dari input user.
func (s *Server) handleExec(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		http.Error(w, "exec failed", http.StatusInternalServerError)
		return
	}
	w.Write(out)
}
