package commands

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/lib"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func HandleXRANGE(conn net.Conn, list_key string, start int64, stop int64) {
	res, err := store.Cache.XRange(list_key, start, stop)
	if err != nil {
		conn.Write([]byte("*0\r\n")) // Return empty array if not found
		return
	}

	var outerArray []any
	for _, entry := range res {
		// Each entry is a [ID, [field1, val1, field2, val2]]
		innerEntry := []any{
			entry.ID,
			flattenFields(entry.Fields),
		}
		outerArray = append(outerArray, innerEntry)
	}

	// This needs to handle nested slices recursively
	conn.Write([]byte(lib.MarshalArrayRESP(outerArray)))
}

// $15\r\n1526985054069-0\r\n

func flattenFields(fields []store.Field) []any {
	var flattened []any
	for _, field := range fields {
		flattened = append(flattened, field.Key, field.Value)
	}
	fmt.Println("Flattened \n", flattened)

	return flattened
}
