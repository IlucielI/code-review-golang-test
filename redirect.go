package main

import (
	"net/http"
)

// Vulnerable: Open Redirect via unvalidated user parameter
func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	http.Redirect(w, r, target, http.StatusFound)
}
