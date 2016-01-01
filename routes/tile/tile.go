package tile

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/log"
	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// Route represents the HTTP route for the resource.
	Route = "/:" + routes.TileType +
		"/:" + routes.TileEndpoint +
		"/:" + routes.TileIndex +
		"/:" + routes.StoreType +
		"/:" + routes.StoreEndpoint +
		"/:" + routes.TileZ +
		"/:" + routes.TileX +
		"/:" + routes.TileY
)

func handleTileErr(w http.ResponseWriter) {
	// send error
	w.WriteHeader(500)
	fmt.Fprint(w, `{"status": "error"}`)
}

// Handler represents the HTTP route response handler.
func Handler(c web.C, w http.ResponseWriter, r *http.Request) {
	// set content type response header
	w.Header().Set("Content-Type", "application/json")
	// parse tile req from URL
	tileReq, err := routes.NewTileRequest(c.URLParams, r.URL.Query())
	if err != nil {
		log.Warn(err)
		handleTileErr(w)
		return
	}
	// parse store req from URL
	storeReq, err := routes.NewStoreRequest(c.URLParams)
	if err != nil {
		log.Warn(err)
		handleTileErr(w)
		return
	}
	// get tile hash
	tileData, err := tile.GetTileFromStore(tileReq, storeReq)
	if err != nil {
		log.Warn(err)
		handleTileErr(w)
		return
	}
	// send response
	w.WriteHeader(200)
	w.Write(tileData)
}
