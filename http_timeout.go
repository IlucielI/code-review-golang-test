package main

import (
	"io"
	"net/http"
)

// Vulnerable: HTTP client without timeout - can leak goroutines and memory
func fetchDataNoTimeout(url string) ([]byte, error) {
	// Default client has no timeout - connections can hang forever
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Vulnerable: Custom client without timeout
var unsafeClient = &http.Client{
	// No Timeout set - connections hang indefinitely
}

func fetchWithUnsafeClient(url string) (string, error) {
	resp, err := unsafeClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

// Vulnerable: Used in loop - compounds the leak
func fetchMultipleURLs(urls []string) []string {
	results := []string{}
	for _, url := range urls {
		// Each hanging request leaks a goroutine
		data, err := fetchDataNoTimeout(url)
		if err != nil {
			continue
		}
		results = append(results, string(data))
	}
	return results
}
