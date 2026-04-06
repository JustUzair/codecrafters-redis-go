package commands

import (
	"fmt"
	"net"
)

func HandlePING(conn net.Conn) {
	val := "+PONG\r\n"
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))

}
