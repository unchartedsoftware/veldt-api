package hbase

import (
	"fmt"
	"math"
	"net/http"

	hbase "github.com/lazyshot/go-hbase"
	log "github.com/unchartedsoftware/plog"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/zenazn/goji/web"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// Route represents the HTTP route for the resource.
	Route = "/hbase" +
		"/:" + TableName +
		"/:" + ZooKeeperHost +
		"/:" + ZooKeeperRoot +
		"/:" + routes.TileZ +
		"/:" + routes.TileX +
		"/:" + routes.TileY
	// column for the tile data blob
	columnName = "tileData:"
)

func handleTileErr(w http.ResponseWriter) {
	// send error
	w.WriteHeader(500)
	fmt.Fprint(w, `{"status": "error"}`)
}

func coordToRowKey(coord *binning.TileCoord) []byte {
	// use the minimum possible number of digits for the tile key
	res := math.Pow(2, float64(coord.Z))
	numDigits := uint32(math.Log10(res) + 1)
	format := fmt.Sprintf("%%02d,%%0%dd,%%0%dd", numDigits, numDigits)
	key := fmt.Sprintf(format, coord.Z, coord.X, coord.Y)
	return []byte(key)
}

// Handler represents the HTTP route response handler.
func Handler(c web.C, w http.ResponseWriter, r *http.Request) {
	// set content type response header
	w.Header().Set("Content-Type", "application/json")
	// parse tile req from URL
	req, err := NewRequest(c.URLParams, r.URL.Query())
	if err != nil {
		log.Warn(err)
		handleTileErr(w)
		return
	}
	// get an hbase client
	client, err := GetClient(req.Host, req.Root)
	if err != nil {
		log.Warn(err)
		handleTileErr(w)
		return
	}
	// get the row key
	row := coordToRowKey(req.Coord)
	// create the get request
	get := hbase.CreateNewGet(row)
	// get data from habse
	result, err := client.Get(req.Table, get)
	if err != nil {
		log.Warn(err)
		handleTileErr(w)
		return
	}
	// extract column
	data, ok := result.Columns[columnName]
	if !ok {
		log.Warnf("Column '%s' is missing from row", columnName)
		handleTileErr(w)
		return
	}
	// send response
	w.WriteHeader(200)
	w.Write(data.Value)
}
