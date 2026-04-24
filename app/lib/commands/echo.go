package commands

import (
	"fmt"
	"io"
)

func HandleECHO(conn io.Writer, val string) {
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
}
