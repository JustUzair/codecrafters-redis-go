package commands

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleBLPOP(conn net.Conn, list_key string, timeout float64) {
	elements, err := store.Cache.BLPop(list_key, timeout)
	if err != nil || elements == nil {
		conn.Write([]byte("*-1\r\n")) // Redis Null Array
		return
	}

	response := lib.MarshalArrayRESP(elements)
	conn.Write([]byte(response))

}
