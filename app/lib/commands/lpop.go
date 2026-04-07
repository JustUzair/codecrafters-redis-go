package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleLPOP(conn net.Conn, list_key string) {
	temp, err := store.Cache.LPop(list_key)
	var element string
	element = temp.(string)
	if err != nil {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	response := fmt.Sprintf("$%d\r\n%s\r\n", len(element), element)
	conn.Write([]byte(response))
}
