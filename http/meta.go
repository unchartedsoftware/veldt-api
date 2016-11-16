package http

import (
	"fmt"
	"net/http"

	"github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism"
	"github.com/zenazn/goji/web"
)

const (
	metaRoute = "/meta/%s"
)

// MetaRoute represents the HTTP route for the resource.
func MetaRoute(pipeline string) string {
	return fmt.Sprintf(metaRoute, pipeline)
}

// MetaHandler represents the HTTP route response handler.
func MetaHandler(pipeline string) web.HandlerFunc {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		// set content type response header
		w.Header().Set("Content-Type", "application/json")
		// parse meta req from URL
		req, err := parseRequestJSON(r.Body)
		if err != nil {
			log.Warn(err)
			handleErr(w, err)
			return
		}
		// ensure it's generated
		err = prism.GenerateMeta(pipeline, req)
		if err != nil {
			log.Warn(err)
			handleErr(w, err)
			return
		}
		// get meta data from store
		meta, err := prism.GetMetaFromStore(pipeline, req)
		if err != nil {
			log.Warn(err)
			handleErr(w, err)
			return
		}
		// send response
		w.WriteHeader(200)
		w.Write(meta)
	}
}
