package meta

import (
	"fmt"
	"net/http"

	log "github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// Route represents the HTTP route for the resource.
	Route = "/meta" +
		"/:" + routes.MetaType +
		"/:" + routes.MetaURI +
		"/:" + routes.StoreType
)

func handleMetaErr(w http.ResponseWriter) {
	// send error
	w.WriteHeader(500)
	fmt.Fprint(w, `{"status": "error"}`)
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
	// ensure it's generated
	err = meta.GenerateMeta(metaReq)
	if err != nil {
		log.Warn(err)
		handleMetaErr(w)
		return
	}
	// get meta data from store
	metaData, err := meta.GetMetaFromStore(metaReq)
	if err != nil {
		log.Warn(err)
		handleMetaErr(w)
		return
	}
	// send response
	w.WriteHeader(200)
	w.Write(metaData)
}
