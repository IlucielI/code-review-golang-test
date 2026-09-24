package main

import (
	"net/http"
	"github.com/gorilla/websocket"
)

// Vulnerable: CheckOrigin returns true blindly allowing Cross-Site WebSocket Hijacking (CSWSH)
var unsafeUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // BUG: CSWSH vulnerability allowing arbitrary cross-origin sites to establish websocket connection
	},
}
