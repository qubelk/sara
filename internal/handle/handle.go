package handle

import (
	"bufio"
	"fmt"
	"net"
	"sara/internal/storage"
	"strconv"
	"strings"
	"time"
)

func Handler(conn net.Conn, store *storage.Storage) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
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
				writer.WriteString("Err: SET command must have at least 2 args, like SET key value\n")
				writer.Flush()
				break
			}

			var ttl time.Duration
			if len(args) >= 4 && strings.ToUpper(args[2]) == "EX" {
				seconds, err := strconv.Atoi(args[3])
				if err != nil {
					writer.WriteString("Err: Expect number of seconds after EX\n")
					writer.Flush()
					break
				}

				if seconds <= 0 {
					writer.WriteString("Err: Time to Life cannot be negative or zero\n")
					writer.Flush()
					break
				}

				ttl = time.Duration(seconds) * time.Second
			} else if len(args) >= 3 {
				writer.WriteString("Err: Expect 2 args for SET command, without ttl specifier\n")
				writer.Flush()
				break
			}

			store.Set(args[0], args[1], ttl)
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
				if v.Value == "" {
					fmt.Fprintf(writer, "%s: (none)\n", k)
				} else {
					fmt.Fprintf(writer, "%s: %s\n", k, v.Value)
				}
			}

			writer.Flush()

		case "DEL":
			if len(args) == 0 {
				writer.WriteString("Err: DEL command must have at least 1 arg, like DEL key\n")
				writer.Flush()
				break
			}

			deleted := store.Del(args...)
			for k, ok := range deleted {
				if ok {
					fmt.Fprintf(writer, "%s: Ok\n", k)
				} else {
					fmt.Fprintf(writer, "%s: (none)\n", k)
				}
			}

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

		case "DUMP":
			tmp := store.Dump()
			if len(tmp) == 0 {
				writer.WriteString("(none)\n")
			}

			for k, v := range tmp {
				fmt.Fprintf(writer, "%s: %s\n", k, v.Value)
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

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(writer, "Failed to scan input: %v", err)
		writer.Flush()
	}
}
