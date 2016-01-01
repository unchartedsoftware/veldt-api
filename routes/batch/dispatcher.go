package batch

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/log"
	"github.com/unchartedsoftware/prism/store"

	"github.com/unchartedsoftware/prism-server/routes"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 256 * 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

// TileDispatcher represents a single clients tile dispatcher.
type TileDispatcher struct {
	RespChan  chan *routes.TileBatchResponse
	ErrChan   chan error
	WaitGroup *sync.WaitGroup
	Conn      *websocket.Conn
}

// NewTileDispatcher returns a pointer to a new tile dispatcher object.
func NewTileDispatcher(w http.ResponseWriter, r *http.Request) (*TileDispatcher, error) {
	// open a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	// set the message read limit
	conn.SetReadLimit(maxMessageSize)
	return &TileDispatcher{
		RespChan:  make(chan *routes.TileBatchResponse),
		ErrChan:   make(chan error),
		WaitGroup: new(sync.WaitGroup),
		Conn:      conn,
	}, nil
}

// ListenAndRespond waits on both tile request and responses and handles each until the websocket connection dies.
func (t *TileDispatcher) ListenAndRespond() error {
	go t.listenForRequests()
	go t.listenForResponses()
	return <-t.ErrChan
}

// Close closes the dispatchers internal channels and websocket connection.
func (t *TileDispatcher) Close() {
	// wait to ensure that no more responses are pending
	t.WaitGroup.Wait()
	// close dispatcher channels
	close(t.RespChan)
	close(t.ErrChan)
	// close websocket connection
	t.Conn.Close()
}

func (t *TileDispatcher) listenForResponses() {
	for res := range t.RespChan {
		tileRes := res.Tile
		// log error if there is one
		if tileRes.Error != nil {
			log.Warn(tileRes.Error)
		}
		// write response to websocket
		t.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		err := t.Conn.WriteJSON(res)
		if err != nil {
			t.ErrChan <- err
			break
		}
	}
}

func (t *TileDispatcher) dispatchRequest(tileReq *tile.Request, storeReq *store.Request) {
	// increment pending response wait group to ensure we don't send on
	// a closed channel
	t.WaitGroup.Add(1)
	// get the tile promise
	promise := tile.GetTile(tileReq, storeReq)
	// when the tile is ready
	promise.OnComplete(func(res interface{}) {
		// cast to tile response
		tileRes := res.(*tile.Response)
		// pass to response channel
		t.RespChan <- routes.NewTileBatchResponse(tileRes, storeReq)
		// decrement the pending response wait group
		t.WaitGroup.Done()
	})
}

func (t *TileDispatcher) getRequest() (*tile.Request, *store.Request, error) {
	// wait on read
	_, msg, err := t.Conn.ReadMessage()
	if err != nil {
		return nil, nil, err
	}
	// parse into tile request
	req, err := routes.NewTileBatchRequest(msg)
	if err != nil {
		// parsing errors should not actually return errors or else the
		// connection will be lost
		return nil, nil, nil
	}
	return req.Tile, req.Store, nil
}

func (t *TileDispatcher) listenForRequests() {
	for {
		// wait on tile request
		tileReq, storeReq, err := t.getRequest()
		if err != nil {
			t.ErrChan <- err
			break
		}
		// if no request could be parsed, return failure immediately
		if tileReq == nil || storeReq == nil {
			t.RespChan <- &routes.TileBatchResponse{
				Tile: &tile.Response{
					Success: false,
				},
			}
			continue
		}
		// dispatch the request
		go t.dispatchRequest(tileReq, storeReq)
	}
}
