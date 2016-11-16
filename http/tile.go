package http

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism"
)

const (
	tileRoute = "/tile/%s"
)

// TileRoute represents the HTTP route for the resource.
func TileRoute(pipeline string) string {
	return fmt.Sprintf(tileRoute, pipeline)
}

// TileHandler represents the HTTP route response handler.
func TileHandler(pipeline string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// set content type response header
		w.Header().Set("Content-Type", "application/json")
		// parse tile req from URL and body
		req, err := parseRequestJSON(r.Body)
		if err != nil {
			log.Warn(err)
			handleErr(w, err)
			return
		}
		// ensure it's generated
		err = prism.GenerateTile(pipeline, req)
		if err != nil {
			log.Warn(err)
			handleErr(w, err)
			return
		}
		// get tile data from store
		tile, err := prism.GetTileFromStore(pipeline, req)
		if err != nil {
			log.Warn(err)
			handleErr(w, err)
			return
		}
		// send response
		w.WriteHeader(200)
		w.Write(tile)
	}
}
