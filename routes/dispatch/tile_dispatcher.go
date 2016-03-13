package dispatch

import (
	"fmt"
	"net/http"

	log "github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism/generation/tile"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// TileRoute represents the HTTP route for the resource.
	TileRoute = "/tile-dispatch"
)

func handleTileRequest(d *Dispatcher, msg []byte) {
	// parse the tile request
	tileReq, err := routes.NewTileBatchRequest(msg)
	if err != nil {
		// parsing error, send back a failure response
		err := fmt.Errorf("Unable to parse tile request message: %s", string(msg))
		// log error
		log.Warn(err)
		err = d.SendResponse(&routes.TileResponse{
			Success: false,
			Err:     err,
		})
		if err != nil {
			log.Warn(err)
		}
		return
	}
	// generate tile and wait on response
	err = tile.GenerateTile(tileReq)
	if err != nil {
		log.Warn(err)
	}
	// create response
	res := routes.NewTileResponse(tileReq, err)
	// send response
	err = d.SendResponse(res)
	if err != nil {
		log.Warn(err)
	}
}

// TileHandler represents the HTTP route response handler.
func TileHandler(w http.ResponseWriter, r *http.Request) {
	// create dispatcher
	dispatcher, err := NewDispatcher(w, r, handleTileRequest)
	if err != nil {
		log.Warn(err)
		return
	}
	// listen for requests and respond
	err = dispatcher.ListenAndRespond()
	if err != nil {
		log.Info(err)
	}
	// clean up dispatcher internals
	dispatcher.Close()
}
