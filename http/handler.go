package http

import (
	"github.com/zenazn/goji/web"
	"net/http"
)

// HandlerFunc represents a handler function for an HTTP endpoint.
type HandlerFunc func(c web.C, w http.ResponseWriter, r *http.Request)
