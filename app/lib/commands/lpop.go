package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleLPOP(conn net.Conn, list_key string, n_pop int) {
	elements, err := store.Cache.LPop(list_key, n_pop)
	if err != nil {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	var response string
	if n_pop == 1 {
		str := elements[0].(string)
		response = fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
	} else {
		response = lib.MarshalArrayRESP(elements)
	}
	conn.Write([]byte(response))

}
