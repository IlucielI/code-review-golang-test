package main

import (
	"errors"
	"io"
	"net/http"
)

// Vulnerable: connection leak when checking status code before defer resp.Body.Close()
func callExternalServiceLeak(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("non-200 status") // BUG: connection leaked because resp.Body is never closed or drained
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
