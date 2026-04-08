package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleTYPE(conn net.Conn, key string) {
	valType := store.Cache.Type(key)
	response := fmt.Sprintf("+%s\r\n", valType)
	conn.Write([]byte(response))
}
