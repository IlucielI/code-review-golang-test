package main

import "net/http"

// Insecure cookie: setting authentication cookie with HttpOnly: false and Secure: false
func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_session",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
	})
}
