package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleGET(conn net.Conn, key string) {
	val, err := store.Cache.Get(key)
	if err != nil {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	// Format: $length\r\ndata\r\n
	switch v := val.(type) {
	case string:
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)))
	}

}
