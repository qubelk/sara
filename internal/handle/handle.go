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

		case "HELLO":
			writer.WriteString("HI\n")
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

			values := store.Get(args...)

			for k, v := range values {
				if v == "" {
					fmt.Fprintf(writer, "%s: (none)\n", k)
				} else {
					fmt.Fprintf(writer, "%s: %s\n", k, v)
				}
			}

			writer.Flush()

		case "DEL":
			if len(args) == 0 {
				writer.WriteString("Err: DEL command must have at least 1 arg, like DEL key\n")
				writer.Flush()
				break
			}

			store.Del(args...)

			writer.WriteString("Ok\n")
			writer.Flush()

		case "EXISTS":
			if len(args) == 0 {
				writer.WriteString("Err: EXISTS command must have at least 1 arg, like EXISTS key\n")
				writer.Flush()
				break
			}

			exists := store.Exists(args...)
			for k, v := range exists {
				fmt.Fprintf(writer, "%s: %v\n", k, v)
			}

			writer.Flush()

		case "BYE", "QUIT":
			writer.WriteString("BYE\n")
			writer.Flush()
			conn.Close()

		default:
			fmt.Fprintf(writer, "Err: Unknown command %s\n", fields[0])
			writer.Flush()
		}
	}
}
