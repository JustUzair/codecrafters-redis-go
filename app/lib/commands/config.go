package commands

import (
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleCONFIG_GET(conn io.Writer, option string) {
	temp, err := store.Cache.ConfigGet(option)
	res := make([]any, len(temp))
	for i, v := range temp {
		res[i] = v
	}
	if err != nil {
		conn.Write([]byte("*0\r\n")) // Return empty array if not found
		return
	}
	conn.Write([]byte(lib.MarshalArrayRESP(res)))

}
