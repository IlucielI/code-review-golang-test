package main

import (
    "context"
    "time"
)

// Vulnerable: Goroutine leak via unbuffered channel without receiver on timeout (CWE-400)
func ProcessWithTimeout(ctx context.Context, data string) string {
    ch := make(chan string) // Unbuffered channel

    go func() {
        // Simulating heavy processing
        time.Sleep(500 * time.Millisecond)
        // If context times out, this write blocks forever and goroutine leaks
        ch <- "processed: " + data
    }()

    select {
    case res := <-ch:
        return res
    case <-ctx.Done():
        return "timeout"
    }
}
