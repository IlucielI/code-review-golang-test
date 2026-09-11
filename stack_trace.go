package main

import (
	"fmt"
	"net/http"
)

// Vulnerable: Error handler exposes stack traces to client
func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	result, err := processRequest(r)
	if err != nil {
		// Dangerous: Full error with stack trace sent to client
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(result))
}

// Vulnerable: Panic handler reveals internal paths
func panicHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			// Exposes internal stack trace with file paths
			http.Error(w, fmt.Sprintf("Panic: %v", err), http.StatusInternalServerError)
		}
	}()

	// Intentional panic to demonstrate
	var data map[string]string
	_ = data["key"] // nil pointer panic
}

// Vulnerable: Database error reveals schema
func queryUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")

	// If SQL error occurs, exposes table/column names
	result, err := executeQuery("SELECT * FROM internal_users WHERE id = ?", userID)
	if err != nil {
		// Database error with schema info sent to client
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(result))
}

func processRequest(r *http.Request) (string, error) {
	return "", fmt.Errorf("internal error at /home/user/app/handlers.go:42")
}

func executeQuery(query string, args ...interface{}) (string, error) {
	return "", fmt.Errorf("pq: column 'ssn' does not exist in table 'internal_users'")
}
