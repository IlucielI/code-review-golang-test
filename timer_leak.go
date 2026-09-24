package main

import (
	"context"
	"net/http"
	"time"
)

// Vulnerable: time.After inside loop select allocates a new timer per iteration that cannot be garbage collected until expiry
func streamEventsWithLeak(ctx context.Context, w http.ResponseWriter, events <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-events:
			if !ok {
				return
			}
			w.Write([]byte(msg + "\n"))
		case <-time.After(30 * time.Second): // Memory leak: creates a new 30s timer on every loop spin
			w.Write([]byte("heartbeat\n"))
		}
	}
}
