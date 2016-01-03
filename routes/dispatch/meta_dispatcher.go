package dispatch

import (
	"errors"
	"net/http"

	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/log"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// MetaRoute represents the HTTP route for the resource.
	MetaRoute = "/meta-dispatch"
)

func handleMetaRequest(d *Dispatcher, msg []byte) {
	// parse the meta request
	req, err := routes.NewMetaBatchRequest(msg)
	if err != nil {
		// parsing error, send back a failure response
		err := d.SendResponse(&routes.MetaBatchResponse{
			Success: false,
			Err:     errors.New("Unable to parse message"),
		})
		log.Warn(err)
		return
	}
	// generate meta data and wait on response
	err = meta.GenerateMeta(req.Meta, req.Store)
	if err != nil {
		log.Warn(err)
	}
	// create response
	res := routes.NewMetaBatchResponse(req.Meta, req.Store, err)
	// send response
	err = d.SendResponse(res)
	if err != nil {
		log.Warn(err)
	}
}

// MetaHandler represents the HTTP route response handler.
func MetaHandler(w http.ResponseWriter, r *http.Request) {
	// create dispatcher
	dispatcher, err := NewDispatcher(w, r, handleMetaRequest)
	if err != nil {
		log.Warn(err)
		return
	}
	// listen for requests and respond
	err = dispatcher.ListenAndRespond()
	if err != nil {
		log.Debug(err)
	}
	// clean up dispatcher internals
	dispatcher.Close()
}
