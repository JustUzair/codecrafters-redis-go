package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleLPUSH(conn net.Conn, list_key string, vals []string) {
	var list_size int = store.Cache.Push(list_key, vals, true)
	conn.Write([]byte(fmt.Sprintf(":%d\r\n", list_size)))
}
