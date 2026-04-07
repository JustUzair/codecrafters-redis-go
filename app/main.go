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
			break
		}

		// fmt.Println("Args %v", args)
		// fmt.Println("Err %v", err)
		command := args[0]
		switch command {
		case "PING":
			commands.HandlePING(conn)
		case "ECHO":
			echoValue := args[1]
			commands.HandleECHO(conn, echoValue)
		case "RPUSH":
			list_key := args[1]
			values := args[2:]
			commands.HandleRPUSH(conn, list_key, values)
		case "LRANGE":
			list_key := args[1]
			start, nilStart := strconv.ParseInt(args[2], 10, 64)
			stop, nilStop := strconv.ParseInt(args[3], 10, 64)
			if nilStart != nil || nilStop != nil {
				fmt.Printf("start and stop indexes are required")
				break
			}
			commands.HandleLRANGE(conn, list_key, start, stop)
		case "SET":
			key := args[1]
			value := args[2]
			if len(args) >= 5 {
				var isDeadlineMillis bool
				flag := args[3]
				deadline := args[4]
				if strings.ToUpper(flag) == "PX" {
					isDeadlineMillis = true
				} else if strings.ToUpper(flag) == "MX" {
					isDeadlineMillis = false
				} else {
					fmt.Printf("Invalid Deadline Parameter %s\n", flag)
					break
				}
				expiry, err := strconv.ParseInt(deadline, 10, 64)
				if err != nil {
					fmt.Printf("Error while parsing deadline: %s\n", deadline)
					break
				}
				commands.HandleSET(conn, key, value, expiry, isDeadlineMillis)
			} else {
				commands.HandleSET(conn, key, value, -1, false)

			}
		case "GET":
			key := args[1]
			commands.HandleGET(conn, key)

		}

	}

}
