package dispatch

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 256 * 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Dispatcher represents a single clients tile dispatcher.
type Dispatcher struct {
	conn    *websocket.Conn
	handler RequestHandler
	mutex   sync.Mutex
}

// RequestHandler represents a handler for the ws request.
type RequestHandler func(*Dispatcher, []byte)

// NewDispatcher returns a pointer to a new tile dispatcher object.
func NewDispatcher(w http.ResponseWriter, r *http.Request, handler RequestHandler) (*Dispatcher, error) {
	// open a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	// set the message read limit
	conn.SetReadLimit(maxMessageSize)
	return &Dispatcher{
		conn:    conn,
		handler: handler,
		mutex:   sync.Mutex{},
	}, nil
}

// ListenAndRespond waits on both tile request and responses and handles each
// until the websocket connection dies.
func (d *Dispatcher) ListenAndRespond() error {
	for {
		// wait on read
		_, msg, err := d.conn.ReadMessage()
		if err != nil {
			return err
		}
		// handle the message
		go d.handler(d, msg)
	}
}

// SendResponse will send a json response in a thread safe manner.
func (d *Dispatcher) SendResponse(res interface{}) error {
	// writes are not thread safe
	d.mutex.Lock()
	defer runtime.Gosched()
	defer d.mutex.Unlock()
	// write response to websocket
	d.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return d.conn.WriteJSON(res)
}

// Close closes the dispatchers websocket connection.
func (d *Dispatcher) Close() {
	// close websocket connection
	d.conn.Close()
}
