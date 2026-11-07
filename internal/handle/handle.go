package handle

import (
	"bufio"
	"fmt"
	"net"
	"sara/internal/storage"
	"strings"
)

func Handler(conn net.Conn, store *storage.Storage) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(writer, "Failed to scan input: %v", err)
			writer.Flush()
			continue
		}

		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		cmd := strings.ToUpper(fields[0])
		args := fields[1:]

		switch cmd {
		case "PING":
			writer.WriteString("PONG\n")
			writer.Flush()

		case "SET":
			if len(args) < 2 {
				writer.WriteString("Err: SET command must have 2 args, like SET key value\n")
				writer.Flush()
				break
			}

			store.Set(args[0], args[1])
			writer.WriteString("Ok\n")
			writer.Flush()

		case "GET":
			if len(args) == 0 {
				writer.WriteString("Err: GET command must have at least 1 arg, like GET key\n")
				writer.Flush()
				break
			}

			v := store.Get(args[0])
			if v == "" {
				writer.WriteString("(none)\n")
			} else {
				fmt.Fprintf(writer, "%s\n", v)
			}

			writer.Flush()
		}
	}
}
