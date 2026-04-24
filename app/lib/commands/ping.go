package commands

import "io"

func HandlePING(conn io.Writer) {
	conn.Write([]byte("+PONG\r\n"))
}
