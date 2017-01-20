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

NOTE: Requires [Glide](https://glide.sh) along with [Go](https://golang.org/) version 1.6+.

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

```go
package main

import (
	"github.com/zenazn/goji/web"
	"github.com/unchartedsoftware/prism-server/http"
	"github.com/unchartedsoftware/prism-server/middleware"
	"github.com/unchartedsoftware/prism-server/ws"
)

func main() {

	// Mount logger middleware
	goji.Use(middleware.Log)
	// Mount gzip middleware
	goji.Use(middleware.Gzip)

	// Meta websocket handler
	log.Infof("Meta WebSocket route: '%s'", ws.MetaRoute)
	goji.Get(ws.MetaRoute, ws.MetaHandler)

	// Tile websocket handler
	log.Infof("Tile WebSocket route: '%s'", ws.TileRoute)
	goji.Get(ws.TileRoute, ws.TileHandler)

	// Meta request handler
	log.Infof("Meta HTTP route: '%s'", http.MetaRoute)
	goji.Post(http.MetaRoute, http.MetaHandler)
	// Tile request handler
	log.Infof("Tile HTTP route: '%s'", http.TileRoute)
	goji.Post(http.TileRoute, http.TileHandler)

	// Start the server
	goji.Serve()
}
```
