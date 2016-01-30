package hbase

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/unchartedsoftware/prism/binning"

	"github.com/unchartedsoftware/prism-server/conf"
	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	// TableName represents the hbase table name.
	TableName = "table-name"
	// ZooKeeperHost represents the zookeeper host for hbase.
	ZooKeeperHost = "zk-host"
	// ZooKeeperRoot represents the zookeeper root for hbase.
	ZooKeeperRoot = "zk-root"
)

// Request represents a HBase tile request.
type Request struct {
	Coord *binning.TileCoord
	Table string
	Host  string
	Root  string
}

// NewRequest instantiates a tile request from a URL and query params.
func NewRequest(params map[string]string, queryParams url.Values) (*Request, error) {
	x, ex := strconv.ParseUint(params[routes.TileX], 10, 32)
	y, ey := strconv.ParseUint(params[routes.TileY], 10, 32)
	z, ez := strconv.ParseUint(params[routes.TileZ], 10, 32)
	if ex != nil || ey != nil || ez != nil {
		return nil, errors.New("Unable to parse tile coordinate request")
	}
	table, ok := params[TableName]
	if !ok {
		return nil, errors.New("Table name missing from request")
	}
	host, ok := params[ZooKeeperHost]
	if !ok {
		return nil, errors.New("Zookeeper host missing from request")
	}
	root, ok := params[ZooKeeperRoot]
	if !ok {
		return nil, errors.New("Zookeeper root missing from request")
	}
	return &Request{
		Coord: &binning.TileCoord{
			X: uint32(x),
			Y: uint32(y),
			Z: uint32(z),
		},
		Table: conf.Unalias(table),
		Host:  conf.Unalias(host),
		Root:  conf.Unalias(root),
	}, nil
}
