package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleSET(conn net.Conn, key string, value string) {
	store.Cache.Set(key, value)
	conn.Write([]byte("+OK\r\n"))
}
