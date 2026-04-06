package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleSET(conn net.Conn, key string, value string, expiry int64, isDeadlineMillis bool) {
	store.Cache.Set(key, value, expiry, isDeadlineMillis)
	conn.Write([]byte("+OK\r\n"))
}
