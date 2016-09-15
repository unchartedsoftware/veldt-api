package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"strconv"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/generation/tile"
)

const (
	// TileType represents the tile "type" component of the URL.
	TileType = "tile-type"
	// TileIndex represents the tile "index" component of the URL.
	TileIndex = "tile-index"
	// TileURI represets the tile "uri" component of the URL.
	TileURI = "tile-uri"
	// TileX represents the tile "x" component of the URL.
	TileX = "x"
	// TileY represents the tile "y" component of the URL.
	TileY = "y"
	// TileZ represents the tile "z" component of the URL.
	TileZ = "z"

	// MetaType represents the meta "type" component of the URL.
	MetaType = "meta-type"
	// MetaIndex represents the meta "index" component of the URL.
	MetaIndex = "meta-index"

	// StoreType represents the store "type" component of the URL.
	StoreType = "store-type"
)

// TileResponse represents a batched tile response.
type TileResponse struct {
	tile.Request
	Success bool  `json:"success"`
	Err     error `json:"-"`
}

// MetaResponse represents a batched meta data response.
type MetaResponse struct {
	meta.Request
	Success bool  `json:"success"`
	Err     error `json:"-"`
}

func parseQueryParams(params url.Values) map[string]interface{} {
	p := make(map[string]interface{})
	for k := range params {
		p[k] = params.Get(k)
	}
	return p
}

// NewTileBatchRequest instantiates a batched tile request from a byte slice.
func NewTileBatchRequest(msg []byte) (*tile.Request, error) {
	req := &tile.Request{}
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// NewTileResponse instantiates a batched tile response.
func NewTileResponse(tileReq *tile.Request, err error) *TileResponse {
	r := &TileResponse{}
	r.Type = tileReq.Type
	r.URI = tileReq.URI
	r.Store = tileReq.Store
	r.Coord = tileReq.Coord
	r.Params = tileReq.Params
	r.Success = (err == nil)
	r.Err = err
	return r
}

// NewMetaBatchRequest instantiates a batched meta data request from a byte
// slice.
func NewMetaBatchRequest(msg []byte) (*meta.Request, error) {
	req := &meta.Request{}
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// NewMetaResponse instantiates a batched meta data reponse.
func NewMetaResponse(metaReq *meta.Request, err error) *MetaResponse {
	r := &MetaResponse{}
	r.Type = metaReq.Type
	r.Index = metaReq.Index
	r.Store = metaReq.Store
	r.Success = (err == nil)
	r.Err = err
	return r
}

// NewTileRequest instantiates a tile request from a URL and body.
func NewTileRequest(url map[string]string, body io.ReadCloser) (*tile.Request, error) {
	x, ex := strconv.ParseUint(url[TileX], 10, 32)
	y, ey := strconv.ParseUint(url[TileY], 10, 32)
	z, ez := strconv.ParseUint(url[TileZ], 10, 32)
	if ex != nil || ey != nil || ez != nil {
		return nil, errors.New("Unable to parse coordinate from tile request")
	}
	typ, ok := url[TileType]
	if !ok {
		return nil, errors.New("Type missing from tile request")
	}
	uri, ok := url[TileURI]
	if !ok {
		return nil, errors.New("URI missing from tile request")
	}
	store, ok := url[StoreType]
	if !ok {
		return nil, errors.New("Store missing from tile request")
	}
	// parse params map
	decoder := json.NewDecoder(body)
	params := make(map[string]interface{})
	err := decoder.Decode(&params)
	if err != nil {
		return nil, err
	}
	return &tile.Request{
		Type: typ,
		URI: uri,
		Coord: &binning.TileCoord{
			X: uint32(x),
			Y: uint32(y),
			Z: uint32(z),
		},
		Params: params,
		Store:  store,
	}, nil
}

// NewMetaRequest instantiates a meta data request from a URL and query params.
func NewMetaRequest(params map[string]string) (*meta.Request, error) {
	typ, ok := params[MetaType]
	if !ok {
		return nil, errors.New("Type missing from meta request")
	}
	index, ok := params[MetaIndex]
	if !ok {
		return nil, errors.New("Index missing from meta request")
	}
	store, ok := params[StoreType]
	if !ok {
		return nil, errors.New("Store missing from meta request")
	}
	return &meta.Request{
		Type:  typ,
		Index: index,
		Store: store,
	}, nil
}
