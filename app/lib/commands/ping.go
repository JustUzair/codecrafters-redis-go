package commands

import (
	"net"
)

func HandlePING(conn net.Conn) {
	conn.Write([]byte("+PONG\r\n"))
}
