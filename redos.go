package main

import (
	"regexp"
)

// Vulnerable: Catastrophic backtracking regex - ReDoS
var emailRegex = regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)+@([a-zA-Z0-9_\-\.]+)+\.([a-zA-Z]{2,5})$`)

func validateEmail(email string) bool {
	// Exponential time complexity on malicious input
	// Example: "aaaaaaaaaaaaaaaaaaaaaaaaaa!" causes catastrophic backtracking
	return emailRegex.MatchString(email)
}

// Vulnerable: Nested quantifiers
var urlRegex = regexp.MustCompile(`^(https?://)?([\da-z\.-]+)+\.([a-z\.]{2,6})+([\/\w \.-]*)*\/?$`)

func validateURL(url string) bool {
	// Multiple nested quantifiers - exponential complexity
	return urlRegex.MatchString(url)
}

// Vulnerable: Used in public endpoint without rate limiting
func validateUserInput(input string) bool {
	// Attacker can send crafted input to cause CPU exhaustion
	return emailRegex.MatchString(input) || urlRegex.MatchString(input)
}
