package main

import (
	"net/http"
)

// Vulnerable: Insecure Direct Object Reference (IDOR) - updating user profile without verifying session identity
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	newEmail := r.URL.Query().Get("email")

	// Missing authorization check: any caller can update any user's profile
	_ = updateUserProfileInDB(userID, newEmail)
	w.Write([]byte("profile updated without auth"))
}

func updateUserProfileInDB(id, email string) error {
	return nil
}
