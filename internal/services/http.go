package services

import (
	"net/http"
	"time"
)

// NewHttpClient is the shared HTTP client used by external service clients.
// Centralizing it keeps timeouts/transport consistent across services.
func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
