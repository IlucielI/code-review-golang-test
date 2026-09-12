package main

import "os/exec"

// Command execution via bash -c with unescaped script
func runUserScript(script string) ([]byte, error) {
	return exec.Command("bash", "-c", script).Output()
}
