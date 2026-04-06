package commands

import "net"

func HandleECHO(conn net.Conn, val string) {
	conn.Write([]byte(val))
}
