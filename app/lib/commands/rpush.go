package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleRPUSH(conn net.Conn, list_key string, vals []string) {
	var list_size int
	for _, val := range vals {
		list_size = store.Cache.RPush(list_key, val)
	}
	conn.Write([]byte(fmt.Sprintf(":%d\r\n", list_size)))
}
