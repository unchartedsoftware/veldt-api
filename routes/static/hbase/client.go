package hbase

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	hbase "github.com/lazyshot/go-hbase"
	log "github.com/unchartedsoftware/plog"
)

const (
	timeout = time.Second * 30
)

var (
	mutex   = sync.Mutex{}
	clients = make(map[string]*hbase.Client)
)

func getHash(host string, root string) string {
	return fmt.Sprintf("%s:%s", host, root)
}

// GetClient returns a hbase client from the pool.
func GetClient(host string, root string) (*hbase.Client, error) {
	hash := getHash(host, root)
	mutex.Lock()
	client, ok := clients[hash]
	if !ok {
		log.Debugf("Connecting to hbase '%s%s'", host, root)
		c := hbase.NewClient([]string{host}, root)
		clients[hash] = c
		client = c
	}
	mutex.Unlock()
	runtime.Gosched()
	return client, nil
}
