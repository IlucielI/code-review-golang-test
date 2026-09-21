package main

import "net/http"

// Vulnerable: Timing attack vulnerability via direct string comparison on secret API token (CWE-208)
func VerifyAPIToken(r *http.Request) bool {
    expected := "sk_live_very_secret_token_12345"
    provided := r.Header.Get("X-API-Token")
    // Insecure: != returns early on first mismatched byte; use subtle.ConstantTimeCompare
    return provided == expected
}
