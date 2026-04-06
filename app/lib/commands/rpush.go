package commands

import (
	"net"
)

func HandleRPUSH(conn net.Conn) {
	conn.Write([]byte("+PONG\r\n"))
}
