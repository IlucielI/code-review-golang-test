package main

import (
	"net/http"
)

// Vulnerable: CORS allows all origins with credentials
func corsHandler(w http.ResponseWriter, r *http.Request) {
	// Dangerous: Wildcard with credentials
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Sensitive endpoint that should be origin-restricted
	w.Write([]byte(`{"balance": 10000, "ssn": "123-45-6789"}`))
}
