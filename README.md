# Prism-Server

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified and `$GOPATH/bin` in your `PATH`.

## Installation

### Using `go get`:

If your project does not use the vendoring tool [Glide](https://glide.sh) to manage dependencies, you can install this package like you would any other:

```bash
go get github.com/unchartedsoftware/prism
```

While this is the simplest way to install the package, due to how `go get` resolves transitive dependencies it may result in version incompatibilities.

### Using `glide get`:

This is the recommended way to install the package and ensures all transitive dependencies are resolved to their compatible versions.

```bash
glide get github.com/unchartedsoftware/prism
```

NOTE: Requires [Glide](https://glide.sh) along with [Go](https://golang.org/) version 1.6, or version 1.5 with the `GO15VENDOREXPERIMENT` environment variable set to `1`.

## Development

Clone the repository:

```bash
mkdir $GOPATH/src/github.com/unchartedsoftware
cd $GOPATH/src/github.com/unchartedsoftware
git clone git@github.com:unchartedsoftware/prism-server.git
```

Install dependencies

```bash
cd prism-server
make install
```

## Usage

This package provides a suite of HTTP and WebSocket handlers to connect the custom live tiling analytics of [prism](https://github.com/unchartedsoftware/prism/) to expressive REST and WebSocket endpoints. This package leverages the power of [goji](https://goji.io/), a minimalistic and feature rich web framework.

## Example

This minimalistic application shows how to register custom tile and meta data generators and connect them to a redis store.

```go
package main

import (
	"net/http"

    "github.com/zenazn/goji"

	"github.com/unchartedsoftware/prism/generation/elastic"
	"github.com/unchartedsoftware/prism/generation/meta"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/store"
	"github.com/unchartedsoftware/prism/store/redis"

	"github.com/unchartedsoftware/prism-server/routes/dispatch"
	"github.com/unchartedsoftware/prism-server/routes/meta"
	"github.com/unchartedsoftware/prism-server/routes/tile"
)

func main() {
	// Register the in-memory store
	store.Register("redis", redis.NewConnection("localhost", "6379"))
	// Register a custom tile and meta data generators
	tile.Register("heatmap", elastic.NewHeatmapTile("http://localhost", "9200"))
	meta.Register("default", elastic.NewDefaultMeta("http://localhost", "9200"))
    // Set the dispatching routes, these endpoints are used to initiate tiling
    // and meta data generation requests over websocket, providing full duplex
    // communication and allowing the server to inform the client the moment the
    // data is ready.
    goji.Get(dispatch.MetaRoute, dispatch.MetaHandler)
    goji.Get(dispatch.TileRoute, dispatch.TileHandler)
    // Set the tile request handler, once tile data is ready, this endpoint can
	// be used to get the generated tile data. If no data is ready this endpoint
    // will attempt to generate it.
    goji.Get(tile.Route, tile.Handler)
    // Set the metadata request handler, this will allow the client to request
	// metadata. If no data is ready this endpoint will attempt to generate it.
    goji.Get(meta.Route, meta.Handler)
	// Greedy route last for static serving
	goji.Get("/*", http.FileServer(http.Dir("./public")))
	// Start the server
	goji.Serve()
}
```

Start the server:

```bash
go run main.go
```

Generate meta data:

```bash
curl -X GET 'http://localhost:8080/default/test_index/redis'
```

This HTTP request results in the following actions:
- Generation of meta data using the `default` generator on the `test_index` index of an elasticsearch instance running on `http://localhost:9200`
- Caching of the generated data in a `redis` store running on `localhost:6379`.

Generate a tile:

```bash
curl -X GET 'http://localhost:8080/heatmap/test_index/redis/4/12/12'
```

This HTTP request results in the following actions:
- Generation of a tile using the `heatmap` generator on the `test` index of an elasticsearch instance running on `http://localhost:9200`
- Caching of the generated tile in a `redis` store running on `localhost:6379`.
