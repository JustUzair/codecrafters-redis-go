package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	_ "strings"

	"github.com/codecrafters-io/redis-starter-go/app/lib/commands"
	_ "github.com/codecrafters-io/redis-starter-go/app/lib/commands"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Printf("Server started at %v\n", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		// fmt.Println("TCP Handshake Successful...")
		go handleConn(conn)
	}

	// var buf []byte

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				// fmt.Println("Client closed the connection")
				break
			}
			fmt.Println("Error reading: ", err.Error())
			return
		}

		if n > 0 {
			var str string = string(buf[:n])
			if strings.Contains(str, "PING") {
				commands.HandlePING(conn)
			} else if strings.Contains(str, "ECHO") {
				vals := strings.Split(str, " ")
				commands.HandleECHO(conn, vals[1])
			} else {
				fmt.Println("Unsupported Command")
				os.Exit(1)
			}

		}

	}

}
