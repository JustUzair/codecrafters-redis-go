package lib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Example RESP string: *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
func UnmarshalRESP(reader *bufio.Reader) ([]string, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	} // Exit loop if client disconnects
	if firstByte != '$' && firstByte != '*' {
		// fmt.Printf("Expected '$', but got: %q (ASCII: %v)\n", firstByte, firstByte)
		return nil, fmt.Errorf("Expecting string to start with '$' or '*'")
	}

	if firstByte == '*' { // *
		var args []string
		sizes, _ := reader.ReadBytes('\n')                                      // 2
		arraySize, _ := strconv.Atoi(strings.TrimSuffix(string(sizes), "\r\n")) // consume 2 and discard \r\n

		// Process [$4\r\nECHO\r\n, $4\r\nECHO\r\n]
		for i := 0; i < arraySize; i++ {
			word, err := ReadBulk(reader)
			if err != nil {
				return nil, err
			}
			args = append(args, word)
		}
		return args, nil

	}
	if firstByte == '$' {
		reader.UnreadByte() // put $ back for reading, so bulk read is usable with
		data, err := ReadBulk(reader)
		return []string{data}, err
	}

	return nil, fmt.Errorf("unknown prefix: `%c`", firstByte)

}

// Process $4\r\nECHO\r\n
func ReadBulk(reader *bufio.Reader) (string, error) {
	// consume $ prefix
	prefix, _ := reader.ReadByte()
	// fmt.Println("expected $, got `%c`", prefix)
	if prefix != '$' {
		return "", fmt.Errorf("expected $, got `%c`", prefix)
	}
	sizes, _ := reader.ReadBytes('\n')
	size, _ := strconv.Atoi(strings.TrimSuffix(string(sizes), "\r\n"))
	data := make([]byte, size)
	io.ReadFull(reader, data)
	reader.Discard(2)

	command := string(data)
	return command, nil
}

/*
Example of RESP Array:

*3\r\n
$1\r\n
a\r\n
$1\r\n
b\r\n
$1\r\n
c\r\n
*/

func MarshalArrayRESP(values []any) string {
	res := fmt.Sprintf("*%d\r\n", len(values))

	for _, v := range values {
		switch val := v.(type) {
		case string:
			// Standard Bulk String
			res += fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
		case []any:
			// Recursive call for nested arrays!
			res += MarshalArrayRESP(val)
		case int:
			// Optional: Handle integers if you ever pass them
			res += fmt.Sprintf(":%d\r\n", val)
		default:
			// Fallback for safety
			s := fmt.Sprint(val)
			res += fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
		}
	}
	return res
}

func MarshalErrorRESP(errorString string) string {
	return fmt.Sprintf("-%s\r\n", errorString)
}
