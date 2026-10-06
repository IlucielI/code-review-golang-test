package main

import (
	"net/http"
)

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	newEmail := r.URL.Query().Get("email")

	_ = updateUserProfileInDB(userID, newEmail)
	w.Write([]byte("profile updated without auth"))
}

func updateUserProfileInDB(id, email string) error {
	return nil
}
