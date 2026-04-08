package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleXADD(conn net.Conn, list_key string, id string, fields []store.Field) {
	n_entries, err := store.Cache.XAdd(list_key, id, fields)
	if err != nil || n_entries == 0 {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	var response string
	if n_entries == 1 {
		response = fmt.Sprintf("$%d\r\n%s\r\n", len(id), id)
	} else {
		response = fmt.Sprintf("+%d\r\n", n_entries)
	}
	conn.Write([]byte(response))
}
