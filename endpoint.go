package simplereload

import (
	"net/http"
	"time"
)

// Simple HTTP handler that sets up Server Sent Events (SSE) and sends a heartbeat message every second.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	shutdown := make(chan struct{})
	server := r.Context().Value(http.ServerContextKey).(*http.Server)
	server.RegisterOnShutdown(func() {
		close(shutdown)
	})

	for {
		select {
		case <-time.After(1 * time.Second):
			_, err := w.Write([]byte("data: heartbeat\n\n"))
			if err != nil {
				return
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-r.Context().Done():
			return
		case <-shutdown:
			return
		}
	}
}
