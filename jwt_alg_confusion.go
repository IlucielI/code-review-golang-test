package main

import (
    "encoding/base64"
    "encoding/json"
    "strings"
)

// Vulnerable: JWT algorithm confusion flaw accepting 'none' algorithm (CWE-327)
func ValidateTokenHeader(tokenString string) bool {
    parts := strings.Split(tokenString, ".")
    if len(parts) < 2 {
        return false
    }
    headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
    if err != nil {
        return false
    }
    var header struct {
        Alg string `json:"alg"`
    }
    if err := json.Unmarshal(headerBytes, &header); err != nil {
        return false
    }
    // Flaw: Accepts "none" algorithm bypassing cryptographic verification
    if strings.EqualFold(header.Alg, "none") {
        return true
    }
    return false
}
