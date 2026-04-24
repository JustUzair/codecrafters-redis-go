package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/lib/commands/router"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

var dir *string
var appendOnly *string
var appendDir *string
var appendFile *string
var appendFsync *string

func main() {
	//  ------- AOF Config Flags -------

	var appendonlyValue string

	if store.Cache.Config.Appendonly {
		appendonlyValue = "yes"
	} else {
		appendonlyValue = "no"
	}

	dir = flag.String("dir", string(store.Cache.Config.Dir), "Directory for data storage")
	appendOnly = flag.String("appendonly", appendonlyValue, "Enable/disable appendonly mode")
	appendDir = flag.String("appenddirname", store.Cache.Config.Appenddirname, "Directory name for AOF")
	appendFile = flag.String("appendfilename", store.Cache.Config.Appendfilename, "Filename for AOF")
	appendFsync = flag.String("appendfsync", store.Cache.Config.Appendfsync, "Fsync policy")

	flag.Parse()

	config := store.Config{
		Dir:            *dir,
		Appendonly:     *appendOnly == "yes",
		Appenddirname:  *appendDir,
		Appendfilename: *appendFile,
		Appendfsync:    *appendFsync,
	}
	// fmt.Println("%v", config)
	store.Cache.ConfigSet(
		config,
	)

	// ------- End of AOF Config Flags -------
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
		go HandleConn(conn)
	}

}

func HandleConn(conn net.Conn) {
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

		router.HandleCommand(conn, args)

	}

}
