package ws

import (
	"net/http"
)

// HandlerFunc represents a handler function for an HTTP endpoint.
type HandlerFunc func(w http.ResponseWriter, r *http.Request)
