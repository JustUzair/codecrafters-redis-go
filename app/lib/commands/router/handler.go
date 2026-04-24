package router

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/lib/commands"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleCommand(conn io.Writer, args []string) {

	var rawCommand []any = make([]any, 0)
	command := args[0]
	rawCommand = append(rawCommand, command)

	switch command {
	case "PING":
		commands.HandlePING(conn)
	case "ECHO":
		echoValue := args[1]
		commands.HandleECHO(conn, echoValue)
	case "RPUSH":
		list_key := args[1]
		values := args[2:]
		rawCommand = append(rawCommand, list_key, values)
		store.Cache.WriteToAOF(rawCommand)
		commands.HandleRPUSH(conn, list_key, values)
	case "LPUSH":
		list_key := args[1]
		values := args[2:]
		rawCommand = append(rawCommand, list_key, values)
		store.Cache.WriteToAOF(rawCommand)
		commands.HandleLPUSH(conn, list_key, values)
	case "LRANGE":
		list_key := args[1]
		start, nilStart := strconv.ParseInt(args[2], 10, 64)
		stop, nilStop := strconv.ParseInt(args[3], 10, 64)
		if nilStart != nil || nilStop != nil {
			fmt.Printf("start and stop indexes are required")
			break
		}
		commands.HandleLRANGE(conn, list_key, start, stop)
	case "LLEN":
		list_key := args[1]
		commands.HandleLLEN(conn, list_key)
	case "LPOP":
		list_key := args[1]
		var n_pop int = 1
		var err error
		if len(args) == 3 {
			n_pop, err = strconv.Atoi(args[2])
			if err != nil {
				fmt.Printf("Invalid argument for number of elements to pop")
				break
			}
		}
		rawCommand = append(rawCommand, list_key, n_pop)
		store.Cache.WriteToAOF(rawCommand)
		commands.HandleLPOP(conn, list_key, n_pop)

	case "BLPOP":
		list_key := args[1]
		timeout, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			fmt.Printf("Invalid argument for number of elements to pop")
			break
		}
		rawCommand = append(rawCommand, list_key, timeout)
		store.Cache.WriteToAOF(rawCommand)
		commands.HandleBLPOP(conn, list_key, timeout)
	case "SET":
		key := args[1]
		value := args[2]
		rawCommand = append(rawCommand, args[1], args[2])
		if len(args) >= 5 {
			var isDeadlineMillis bool
			flag := args[3]
			deadline := args[4]
			rawCommand = append(rawCommand, args[3], args[4])

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

			store.Cache.WriteToAOF(rawCommand)
			commands.HandleSET(conn, key, value, expiry, isDeadlineMillis)
		} else {
			store.Cache.WriteToAOF(rawCommand)
			commands.HandleSET(conn, key, value, -1, false)

		}
	case "GET":
		key := args[1]
		commands.HandleGET(conn, key)

	case "TYPE":
		key := args[1]
		commands.HandleTYPE(conn, key)
	case "XADD":
		list_key := args[1]
		stream_id := args[2]
		if len(args) <= 3 {
			fmt.Printf("Insufficient params count, args passed %d\n", len(args))
			return
		}
		var rawFields []string = args[3:]

		if len(rawFields)%2 != 0 {
			fmt.Println("raw fields len", len(rawFields))
			fmt.Println("Every field must consist of a key and a value")
			return
		}
		rawCommand = append(rawCommand, list_key, stream_id, rawFields)
		var fields []store.Field
		var fieldLen int = len(rawFields)

		for i := 0; i < fieldLen-1; i += 2 {
			key := rawFields[i]
			value := rawFields[i+1]
			fields = append(fields, store.Field{Key: string(key), Value: any(value)})
		}
		// fmt.Println("Fields Constructed \n Fields: \n ", fields)
		store.Cache.WriteToAOF(rawCommand)
		commands.HandleXADD(conn, list_key, stream_id, fields)
	case "XRANGE":
		list_key := args[1]
		start := args[2]
		stop := args[3]

		commands.HandleXRANGE(conn, list_key, start, stop)

	case "CONFIG":
		subCommand := args[1]
		option := args[2]
		if subCommand == "GET" {
			// CONFIG GET
			if option == "dir" || option == "appendonly" || option == "appenddirname" || option == "appendfilename" || option == "appendfsync" {
				commands.HandleCONFIG_GET(conn, option)
			} else {
				fmt.Printf("Invalid option: %s\n", option)
			}
		}
	}
}
