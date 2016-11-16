package http

import (
	"fmt"
	"net/http"
)

func handleErr(w http.ResponseWriter, err error) {
	// send error
	w.WriteHeader(500)
	fmt.Fprint(w, `{"status": "error", "error": "`+err.Error()+`"}`)
}
