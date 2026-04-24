package commands

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleLRANGE(conn io.Writer, list_key string, start int64, stop int64) {
	res := store.Cache.LRange(list_key, start, stop)
	respData := lib.MarshalArrayRESP(res)
	conn.Write([]byte(fmt.Sprint(respData)))
}
