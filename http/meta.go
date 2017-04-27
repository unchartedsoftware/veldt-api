package http

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/util/json"
)

const (
	// MetaRoute represents the HTTP route for the resource.
	MetaRoute = "/meta"
)

// MetaHandler represents the HTTP route response handler.
func MetaHandler(w http.ResponseWriter, r *http.Request) {
	// parse meta req from URL and body
	req, err := parseRequestJSON(r.Body)
	if err != nil {
		log.Warn(err)
		handleErr(w, err)
		return
	}
	// get pipeline id
	pipeline, ok := json.GetString(req, "pipeline")
	if !ok {
		// send error response
		err := fmt.Errorf(`no "pipeline" argument is provided`)
		log.Warn(err)
		handleErr(w, err)
		return
	}
	// get meta data from store, ensuring it is generated
	meta, err := veldt.GenerateAndGetMeta(pipeline, req)
	if err != nil {
		log.Warn(err)
		handleErr(w, err)
		return
	}
	// write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(meta)
}
