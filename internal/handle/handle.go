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

func handlePing(writer *bufio.Writer) {
	writer.WriteString("PONG\n")
	writer.Flush()
}

func handleHello(writer *bufio.Writer) {
	writer.WriteString("HI\n")
	writer.Flush()
}

func handleSet(writer *bufio.Writer, store *storage.Storage, args []string) {
	defer writer.Flush()
	if len(args) < 2 {
		writer.WriteString("Err: SET command must have at least 2 args, like SET key value\n")
		return
	}

	var ttl time.Duration
	if len(args) >= 4 && strings.ToUpper(args[2]) == "EX" {
		seconds, err := strconv.Atoi(args[3])
		if err != nil {
			writer.WriteString("Err: Expect number of seconds after EX\n")
			return
		}

		if seconds <= 0 {
			writer.WriteString("Err: Time to Life cannot be negative or zero\n")
			return
		}

		ttl = time.Duration(seconds) * time.Second
	} else if len(args) >= 3 {
		writer.WriteString("Err: Expect 2 args for SET command, without ttl specifier\n")
		return
	}

	store.Set(args[0], args[1], ttl)
	writer.WriteString("Ok\n")
}

func handleGet(writer *bufio.Writer, store *storage.Storage, args []string) {
	defer writer.Flush()
	if len(args) == 0 {
		writer.WriteString("Err: GET command must have at least 1 arg, like GET key\n")
		return
	}

	values := store.Get(args...)

	for k, v := range values {
		if v.Value == "" {
			fmt.Fprintf(writer, "%s: (none)\n", k)
		} else {
			fmt.Fprintf(writer, "%s: %s\n", k, v.Value)
		}
	}
}

func handleDel(writer *bufio.Writer, store *storage.Storage, args []string) {
	defer writer.Flush()
	if len(args) == 0 {
		writer.WriteString("Err: DEL command must have at least 1 arg, like DEL key\n")
		return
	}

	deleted := store.Del(args...)
	for k, ok := range deleted {
		if ok {
			fmt.Fprintf(writer, "%s: Ok\n", k)
		} else {
			fmt.Fprintf(writer, "%s: (none)\n", k)
		}
	}
}

func handleExists(writer *bufio.Writer, store *storage.Storage, args []string) {
	defer writer.Flush()
	if len(args) == 0 {
		writer.WriteString("Err: EXISTS command must have at least 1 arg, like EXISTS key\n")
		return
	}

	exists := store.Exists(args...)
	for k, v := range exists {
		fmt.Fprintf(writer, "%s: %v\n", k, v)
	}
}

func handleDump(writer *bufio.Writer, store *storage.Storage) {
	defer writer.Flush()

	dump := store.Dump()
	if len(dump) == 0 {
		writer.WriteString("(none)\n")
		return
	}

	for k, v := range dump {
		fmt.Fprintf(writer, "%s: %s\n", k, v.Value)
	}
}

func handleBye(writer *bufio.Writer) {
	writer.WriteString("BYE\n")
	writer.Flush()
}

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
			handlePing(writer)

		case "HELLO":
			handleHello(writer)

		case "SET":
			handleSet(writer, store, args)

		case "GET":
			handleGet(writer, store, args)

		case "DEL":
			handleDel(writer, store, args)

		case "EXISTS":
			handleExists(writer, store, args)

		case "DUMP":
			handleDump(writer, store)

		case "BYE", "QUIT":
			handleBye(writer)
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
