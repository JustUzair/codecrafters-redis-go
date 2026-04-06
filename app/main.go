package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/lib/commands"
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

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buffReader := bufio.NewReader(conn)
	for {

		args, err := lib.UnmarshalRESP(buffReader)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			return
		}

		// fmt.Println("Args %v", args)
		// fmt.Println("Err %v", err)
		command := args[0]
		switch command {
		case "PING":
			commands.HandlePING(conn)
		case "ECHO":
			commands.HandleECHO(conn, args[1])
		case "RPUSH":
			commands.HandleRPUSH(conn, args[1], args[2:])
		case "SET":
			if len(args) >= 5 {
				var isDeadlineMillis bool

				if strings.ToUpper(args[3]) == "PX" {
					isDeadlineMillis = true
				} else if strings.ToUpper(args[3]) == "MX" {
					isDeadlineMillis = false
				} else {
					fmt.Printf("Invalid Deadline Parameter %s\n", args[3])
					break
				}
				expiry, err := strconv.ParseInt(args[4], 10, 64)
				if err != nil {
					fmt.Printf("Error while parsing deadline: %s\n", args[4])
					break
				}
				commands.HandleSET(conn, args[1], args[2], expiry, isDeadlineMillis)
			} else {
				commands.HandleSET(conn, args[1], args[2], -1, false)

			}
		case "GET":
			commands.HandleGET(conn, args[1])

		}

	}

}
