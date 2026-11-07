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
		cmd := strings.ToUpper(fields[0])

		switch cmd {
		case "PING":
			writer.WriteString("PONG\n")
			writer.Flush()
		}
	}
}
