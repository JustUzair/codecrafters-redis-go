package commands

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleRPUSH(conn io.Writer, list_key string, value []string) {
	var list_size int = store.Cache.Push(list_key, value, false)
	conn.Write([]byte(fmt.Sprintf(":%d\r\n", list_size)))
}
