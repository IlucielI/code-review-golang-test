package main

import (
	"encoding/json"
	"net/http"
)

// Vulnerable: Public API endpoint without rate limiting
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// No rate limiting - allows brute force attacks
	json.NewDecoder(r.Body).Decode(&creds)

	// Expensive operation without throttling
	if authenticateUser(creds.Username, creds.Password) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": "..."})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// Vulnerable: Password reset without rate limiting
func handlePasswordReset(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	// No rate limiting - allows enumeration and spam
	sendPasswordResetEmail(email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reset email sent"})
}

// Vulnerable: Expensive operation without rate limiting
func handleSearchAPI(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	// No rate limiting - allows DoS via expensive searches
	results := performExpensiveSearch(query)

	json.NewEncoder(w).Encode(results)
}

func authenticateUser(username, password string) bool {
	// Stub
	return username == "admin" && password == "secret"
}

func sendPasswordResetEmail(email string) {}
func performExpensiveSearch(query string) []string { return nil }
