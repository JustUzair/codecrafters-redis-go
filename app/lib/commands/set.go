package commands

import (
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleSET(conn io.Writer, key string, value any, expiry int64, isDeadlineMillis bool) {
	store.Cache.Set(key, value, expiry, isDeadlineMillis)
	conn.Write([]byte("+OK\r\n"))
}
