package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleLLEN(conn net.Conn, list_key string) {
	var list_size int = store.Cache.LLen(list_key)
	conn.Write([]byte(fmt.Sprintf(":%d\r\n", list_size)))
}
