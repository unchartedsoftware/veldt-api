package meta

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/log"
	"github.com/unchartedsoftware/prism/store"
	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// Route represents the HTTP route for the resource.
	Route = "/:" + routes.MetaType +
		"/:" + routes.MetaEndpoint +
		"/:" + routes.MetaIndex +
		"/:" + routes.StoreType +
		"/:" + routes.StoreEndpoint
)

func handleMetaErr(w http.ResponseWriter) {
	// send error
	w.WriteHeader(500)
	fmt.Fprint(w, `{"status": "error"}`)
}

func dispatchRequest(metaChan chan *meta.Response, metaReq *meta.Request, storeReq *store.Request) {
	// get the meta data promise
	promise := meta.GetMeta(metaReq, storeReq)
	// when the meta data is ready
	promise.OnComplete(func(res interface{}) {
		metaChan <- res.(*meta.Response)
	})
}

// Handler represents the HTTP route response handler.
func Handler(c web.C, w http.ResponseWriter, r *http.Request) {
	// set content type response header
	w.Header().Set("Content-Type", "application/json")
	// parse meta req from URL
	metaReq, err := routes.NewMetaRequest(c.URLParams)
	if err != nil {
		log.Warn(err)
		handleMetaErr(w)
		return
	}
	// parse store req from URL
	storeReq, err := routes.NewStoreRequest(c.URLParams)
	if err != nil {
		log.Warn(err)
		handleMetaErr(w)
		return
	}
	// create channel to pass metadata
	metaChan := make(chan *meta.Response)
	// dispatch the request async and wait on channel
	go dispatchRequest(metaChan, metaReq, storeReq)
	// wait on response
	metaRes := <-metaChan
	if metaRes.Error != nil {
		log.Warn(metaRes.Error)
		handleMetaErr(w)
		return
	}
	// send success response
	w.WriteHeader(200)
	w.Write(metaRes.Meta)
}
