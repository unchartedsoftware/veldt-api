package routes

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/store"

	"github.com/unchartedsoftware/prism-server/conf"
)

const (
	// TileType represents the tile "type" component of the URL.
	TileType = "tile-type"
	// TileEndpoint represents the tile "endpoint" component of the URL.
	TileEndpoint = "tile-endpoint"
	// TileIndex represents the tile "index" component of the URL.
	TileIndex = "tile-index"
	// TileX represents the tile "x" component of the URL.
	TileX = "x"
	// TileY represents the tile "y" component of the URL.
	TileY = "y"
	// TileZ represents the tile "z" component of the URL.
	TileZ = "z"

	// MetaType represents the meta "type" component of the URL.
	MetaType = "meta-type"
	// MetaEndpoint represents the meta "endpoint" component of the URL.
	MetaEndpoint = "meta-endpoint"
	// MetaIndex represents the meta "index" component of the URL.
	MetaIndex = "meta-index"

	// StoreType represents the store "type" component of the URL.
	StoreType = "store-type"
	// StoreEndpoint represents the store "endpoint" component of the URL.
	StoreEndpoint = "store-endpoint"
)

// TileBatchRequest represents a batched tile request.
type TileBatchRequest struct {
	Tile  *tile.Request  `json:"tile"`
	Store *store.Request `json:"store"`
}

// TileBatchResponse represents a batched tile response.
type TileBatchResponse struct {
	Tile    *tile.Request  `json:"tile"`
	Store   *store.Request `json:"store"`
	Success bool           `json:"success"`
	Err     error          `json:"-"`
}

// MetaBatchRequest represents a batched meta data request.
type MetaBatchRequest struct {
	Meta  *meta.Request  `json:"meta"`
	Store *store.Request `json:"store"`
}

// MetaBatchResponse  represents a batched meta data response.
type MetaBatchResponse struct {
	Meta    *meta.Request  `json:"meta"`
	Store   *store.Request `json:"store"`
	Success bool           `json:"success"`
	Err     error          `json:"-"`
}

func parseQueryParams(params url.Values) map[string]interface{} {
	p := make(map[string]interface{})
	for k := range params {
		p[k] = params.Get(k)
	}
	return p
}

// NewTileBatchRequest instantiates a batched tile request from a byte slice.
func NewTileBatchRequest(msg []byte) (*TileBatchRequest, error) {
	req := &TileBatchRequest{}
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, err
	}
	// un-alias
	// tile
	req.Tile.Type = conf.Unalias(req.Tile.Type)
	req.Tile.Endpoint = conf.Unalias(req.Tile.Endpoint)
	req.Tile.Index = conf.Unalias(req.Tile.Index)
	// store
	req.Store.Type = conf.Unalias(req.Store.Type)
	req.Store.Endpoint = conf.Unalias(req.Store.Endpoint)
	return req, nil
}

// NewTileBatchResponse instantiates a batched tile response.
func NewTileBatchResponse(tileReq *tile.Request, storeReq *store.Request, err error) *TileBatchResponse {
	req := &TileBatchResponse{
		Tile:    tileReq,
		Store:   storeReq,
		Success: (err == nil),
		Err:     err,
	}
	// alias
	// tile
	req.Tile.Type = conf.Alias(req.Tile.Type)
	req.Tile.Endpoint = conf.Alias(req.Tile.Endpoint)
	req.Tile.Index = conf.Alias(req.Tile.Index)
	// store
	req.Store.Type = conf.Alias(req.Store.Type)
	req.Store.Endpoint = conf.Alias(req.Store.Endpoint)
	return req
}

// NewMetaBatchRequest instantiates a batched meta data request from a byte
// slice.
func NewMetaBatchRequest(msg []byte) (*MetaBatchRequest, error) {
	req := &MetaBatchRequest{}
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, err
	}
	// un-alias
	// meta
	req.Meta.Type = conf.Unalias(req.Meta.Type)
	req.Meta.Endpoint = conf.Unalias(req.Meta.Endpoint)
	req.Meta.Index = conf.Unalias(req.Meta.Index)
	// store
	req.Store.Type = conf.Unalias(req.Store.Type)
	req.Store.Endpoint = conf.Unalias(req.Store.Endpoint)
	return req, nil
}

// NewMetaBatchResponse instantiates a batched meta data reponse.
func NewMetaBatchResponse(metaRes *meta.Request, storeReq *store.Request, err error) *MetaBatchResponse {
	req := &MetaBatchResponse{
		Meta:    metaRes,
		Store:   storeReq,
		Success: (err == nil),
		Err:     err,
	}
	// alias
	// meta
	req.Meta.Type = conf.Alias(req.Meta.Type)
	req.Meta.Endpoint = conf.Alias(req.Meta.Endpoint)
	req.Meta.Index = conf.Alias(req.Meta.Index)
	// store
	req.Store.Type = conf.Alias(req.Store.Type)
	req.Store.Endpoint = conf.Alias(req.Store.Endpoint)
	return req
}

// NewTileRequest instantiates a tile request from a URL and query params.
func NewTileRequest(params map[string]string, queryParams url.Values) (*tile.Request, error) {
	x, ex := strconv.ParseUint(params[TileX], 10, 32)
	y, ey := strconv.ParseUint(params[TileY], 10, 32)
	z, ez := strconv.ParseUint(params[TileZ], 10, 32)
	if ex != nil || ey != nil || ez != nil {
		return nil, errors.New("Unable to parse coordinate from tile request")
	}
	typ, ok := params[TileType]
	if !ok {
		return nil, errors.New("Type missing from tile request")
	}
	endpoint, ok := params[TileEndpoint]
	if !ok {
		return nil, errors.New("Endpoint missing from tile request")
	}
	index, ok := params[TileIndex]
	if !ok {
		return nil, errors.New("Index missing from tile request")
	}
	return &tile.Request{
		Coord: &binning.TileCoord{
			X: uint32(x),
			Y: uint32(y),
			Z: uint32(z),
		},
		Type:     conf.Unalias(typ),
		Endpoint: conf.Unalias(endpoint),
		Index:    conf.Unalias(index),
		Params:   parseQueryParams(queryParams),
	}, nil
}

// NewMetaRequest instantiates a meta data request from a URL and query params.
func NewMetaRequest(params map[string]string) (*meta.Request, error) {
	typ, ok := params[MetaType]
	if !ok {
		return nil, errors.New("Type missing from meta request")
	}
	endpoint, ok := params[MetaEndpoint]
	if !ok {
		return nil, errors.New("Endpoint missing from meta request")
	}
	index, ok := params[MetaIndex]
	if !ok {
		return nil, errors.New("Index missing from meta request")
	}
	return &meta.Request{
		Type:     conf.Unalias(typ),
		Endpoint: conf.Unalias(endpoint),
		Index:    conf.Unalias(index),
	}, nil
}

// NewStoreRequest instantiates a store request from a URL and query params.
func NewStoreRequest(params map[string]string) (*store.Request, error) {
	typ, ok := params[StoreType]
	if !ok {
		return nil, errors.New("Type missing from store request")
	}
	endpoint, ok := params[StoreEndpoint]
	if !ok {
		return nil, errors.New("Endpoint missing from store request")
	}
	return &store.Request{
		Type:     conf.Unalias(typ),
		Endpoint: conf.Unalias(endpoint),
	}, nil
}
