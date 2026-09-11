package main

import (
	"io"
	"net/http"
	"os"
)

// fetchRemote downloads a remote URL to a local file.
func (s *Server) fetchRemote(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	return err
}
