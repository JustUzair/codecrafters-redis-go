package commands

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleXADD(conn io.Writer, list_key string, id string, fields []store.Field) {
	ok, id, err := store.Cache.XAdd(list_key, id, fields)
	if !ok && err != nil {
		conn.Write([]byte(lib.MarshalErrorRESP(err.Error())))
		return
	}
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(id), id)))
}
