package commands

import (
	"fmt"
	"net"
)

func HandleECHO(conn net.Conn, val string) {
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
}
