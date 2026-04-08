package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleXADD(conn net.Conn, list_key string, id string, fields []store.Field) {
	err := store.Cache.XAdd(list_key, id, fields)
	if err != nil {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	response := fmt.Sprintf("$%d\r\n%s\r\n", len(id), id)
	conn.Write([]byte(response))
}
