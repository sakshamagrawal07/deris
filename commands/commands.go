package commands

import (
	"log"
	"net"
	"strings"

	"github.com/sakshamagrawal07/deris/data"
)

func ReadCommand(c net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil

}

func parseCommand(cmd string) []string {
	log.Println("Parsing command...")
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	return strings.Fields(cmd)
}

func ExecuteCommand(cmd string) string {
	log.Println("Executing command")
	parsedCommand := parseCommand(cmd)
	log.Println("Parsed command: ", parsedCommand)
	cmd = parsedCommand[0]

	var response string
	switch cmd {
	case "set":
		data.Set(parsedCommand[1], parsedCommand[2])
		response = "OK\n"
	case "get":
		value, ok := data.Get(parsedCommand[1])
		if ok {
			response = "\"" + value + "\"" + "\n"
		} else {
			response = "(nil)\n"
		}
	case "exit":
		response = "Bye\n"
	default:
		response = "ERR unknown command '" + cmd + "'"
	}
	return response
}
