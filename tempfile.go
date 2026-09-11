package main

import (
	"net/http"
	"os"
	"path/filepath"
)

// Vulnerable: Path traversal via unsafe file creation from user input
func HandleCreateTempFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	targetPath := filepath.Join("/tmp/uploads", filename)

	f, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	w.Write([]byte("created"))
}
