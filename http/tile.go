package http

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/util/json"
)

const (
	// TileRoute represents the HTTP route for the resource.
	TileRoute = "/tile"
)

// TileHandler represents the HTTP route response handler.
func TileHandler(w http.ResponseWriter, r *http.Request) {
	// parse tile req from URL and body
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
	// get tile data from store, ensuring it is generated
	tile, err := veldt.GenerateAndGetTile(pipeline, req)
	if err != nil {
		log.Warn(err)
		handleErr(w, err)
		return
	}
	// write response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(tile)
}
