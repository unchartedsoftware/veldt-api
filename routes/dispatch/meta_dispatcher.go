package dispatch

import (
	"fmt"
	"net/http"

	log "github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism/generation/meta"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// MetaRoute represents the HTTP route for the resource.
	MetaRoute = "/meta-dispatch"
)

func handleMetaRequest(d *Dispatcher, msg []byte) {
	// parse the meta request
	metaReq, err := routes.NewMetaBatchRequest(msg)
	if err != nil {
		// parsing error, send back a failure response
		err := fmt.Errorf("Unable to parse meta request message: %s", string(msg))
		// log error
		log.Warn(err)
		err = d.SendResponse(&routes.MetaResponse{
			Success: false,
			Err:     err,
		})
		if err != nil {
			log.Warn(err)
		}
		return
	}
	// generate meta data and wait on response
	err = meta.GenerateMeta(metaReq)
	if err != nil {
		log.Warn(err)
	}
	// create response
	res := routes.NewMetaResponse(metaReq, err)
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
		log.Info(err)
	}
	// clean up dispatcher internals
	dispatcher.Close()
}
