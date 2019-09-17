# veldt-api

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/unchartedsoftware/veldt-api)
[![Build Status](https://travis-ci.org/unchartedsoftware/veldt-api.svg?branch=master)](https://travis-ci.org/unchartedsoftware/veldt-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/unchartedsoftware/veldt-api)](https://goreportcard.com/report/github.com/unchartedsoftware/veldt-api)

## Dependencies

Requires the [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified and `$GOPATH/bin` in your `PATH`.

## Installation

```bash
go get github.com/unchartedsoftware/veldt-api
```

## Development

NOTE: Requires [Go](https://golang.org/) version 1.12+.

Clone the repository outside of your `$GOPATH`:

```bash
git clone git@github.com:unchartedsoftware/veldt-api.git
cd veldt-api
```

Install dependencies

```bash
make install
```

## Usage

This package provides HTTP and WebSocket handlers to connect the on-demand tile-based analytics of [veldt](https://github.com/unchartedsoftware/veldt/) to  HTTP and WebSocket endpoints.

## Example

```go
package main

import (
	"github.com/zenazn/goji"
	"github.com/unchartedsoftware/veldt-api/http"
	"github.com/unchartedsoftware/veldt-api/middleware"
	"github.com/unchartedsoftware/veldt-api/ws"
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
