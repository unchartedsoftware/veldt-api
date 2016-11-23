package ws

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/util/json"
)

const (
	// MetaRoute represents the HTTP route for the resource.
	MetaRoute = "/ws/meta"
)

// MetaHandler represents the HTTP route response handler.
func MetaHandler(w http.ResponseWriter, r *http.Request) {
	// create conn
	conn, err := NewConnection(w, r, handleMetaRequest)
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

func handleMetaRequest(conn *Connection, msg []byte) {
	// parse the meta request into JSON
	req, err := parseRequestJSON(msg)
	if err != nil {
		// parsing error, send back a failure response
		err := fmt.Errorf("unable to parse meta request message: %s", string(msg))
		// log error
		log.Warn(err)
		// send error response
		err = handleErr(conn, err)
		if err != nil {
			log.Warn(err)
		}
		return
	}
	// get pipeline id
	pipeline, ok := json.GetString(req, "pipeline")
	if !ok {
		// send error response
		err = handleErr(conn, err)
		if err != nil {
			log.Warn(err)
		}
		return
	}
	// generate meta data and wait on response
	err = prism.GenerateMeta(pipeline, req)
	if err != nil {
		log.Warn(err)
	}
	// create response by appending success / error fields
	req["success"] = err == nil
	req["error"] = err
	// send response
	err = conn.SendResponse(req)
	if err != nil {
		log.Warn(err)
	}
}
