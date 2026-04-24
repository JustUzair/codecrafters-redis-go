package commands

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleTYPE(conn io.Writer, key string) {
	valType := store.Cache.Type(key)
	response := fmt.Sprintf("+%s\r\n", valType)
	conn.Write([]byte(response))
}
