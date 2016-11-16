package ws

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism"
)

const (
	tileRoute = "/ws/tile/%s"
)

// TileRoute represents the HTTP route for the resource.
func TileRoute(pipeline string) string {
	return fmt.Sprintf(tileRoute, pipeline)
}

// TileHandler represents the HTTP route response handler.
func TileHandler(pipeline string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create conn
		conn, err := NewConnection(w, r, handleTileRequest(pipeline))
		if err != nil {
			log.Warn(err)
			return
		}
		// listen for requests and respond
		err = conn.ListenAndRespond()
		if err != nil {
			log.Info(err)
		}
		// clean up conn internals
		conn.Close()
	}
}

func handleTileRequest(pipeline string) requestHandler {
	return func(conn *Connection, msg []byte) {
		// parse the tile request into JSON
		req, err := parseRequestJSON(msg)
		if err != nil {
			// parsing error, send back a failure response
			err := fmt.Errorf("unable to parse tile request message: %s", string(msg))
			// log error
			log.Warn(err)
			// send error response
			err = conn.SendResponse(map[string]interface{}{
				"success": false,
				"error":   err,
			})
			if err != nil {
				log.Warn(err)
			}
			return
		}
		// generate tile and wait on response
		err = prism.GenerateTile(pipeline, req)
		if err != nil {
			log.Warn(err)
		}
		// create response by appending success / error fields
		req["success"] = err != nil
		req["error"] = err
		// send response
		err = conn.SendResponse(req)
		if err != nil {
			log.Warn(err)
		}
	}
}
