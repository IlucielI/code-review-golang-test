package main

import (
	"context"
	"net/http"
	"time"
)

// Safe: timer is allocated once and properly stopped
func streamEventsSafe(ctx context.Context, w http.ResponseWriter, events <-chan string) {
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-events:
			if !ok {
				return
			}
			w.Write([]byte(msg + "\n"))
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(30 * time.Second)
		case <-timer.C:
			w.Write([]byte("heartbeat\n"))
			timer.Reset(30 * time.Second)
		}
	}
}
