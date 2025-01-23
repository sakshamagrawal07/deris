package commands

import (
	"errors"
	"log"
	"net"
	"strconv"
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

func ExecuteCommand(cmd string) (string, error) {
	log.Println("Executing command")
	parsedCommand := parseCommand(cmd)
	log.Println("Parsed command: ", parsedCommand)

	var response string
	var err error = nil

	if len(parsedCommand) == 0 {
		response = "ERR empty command\n"
		err = errors.New("ERR empty command")
		return response, err
	}

	cmd = parsedCommand[0]

	switch cmd {
	case "set":
		if len(parsedCommand) != 3 {
			response = "ERR wrong number of arguments for 'set' command\nRequired `set <key> <value>`\n"
			err = errors.New("ERR wrong number of arguments for 'set' command\nRequired `set <key> <value>`")
			break
		}
		data.Set(parsedCommand[1], parsedCommand[2])
		response = "OK\n"
	case "setnx":
		if len(parsedCommand) != 3 {
			response = "ERR wrong number of arguments for 'setnx' command\nRequired `setnx <key> <value>`\n"
			err = errors.New("ERR wrong number of arguments for 'setnx' command\nRequired `setnx <key> <value>`")
			break
		}
		data.Setnx(parsedCommand[1], parsedCommand[2])
		response = "OK\n"
	case "get":
		if len(parsedCommand) != 2 {
			response = "ERR wrong number of arguments for 'get' command\nRequired `get <key>`\n"
			err = errors.New("ERR wrong number of arguments for 'get' command\nRequired `get <key>`")
			break
		}
		value, ok := data.Get(parsedCommand[1])
		if ok {
			response = strings.Join(value, " ") + "\n"
		} else {
			response = "(nil)\n"
		}
	case "lpush":
		if len(parsedCommand) != 3 {
			response = "ERR wrong number of arguments for 'lpush' command\nRequired `lpush <key> <value>`\n"
			err = errors.New("ERR wrong number of arguments for 'lpush' command\nRequired `lpush <key> <value>`")
			break
		}
		data.LPush(parsedCommand[1], parsedCommand[2])
		response = "OK\n"
	case "lpop":
		if len(parsedCommand) != 2 {
			response = "ERR wrong number of arguments for 'lpop' command\nRequired `lpop <key>`\n"
			err = errors.New("ERR wrong number of arguments for 'lpop' command\nRequired `lpop <key>`")
			break
		}
		value, ok := data.LPop(parsedCommand[1])
		if ok {
			response = "\"" + value + "\"" + "\n"
		} else {
			response = "(nil)\n"
		}
	case "rpush":
		if len(parsedCommand) != 3 {
			response = "ERR wrong number of arguments for 'rpush' command\nRequired `rpush <key> <value>`\n"
			err = errors.New("ERR wrong number of arguments for 'rpush' command\nRequired `rpush <key> <value>`")
			break
		}
		data.RPush(parsedCommand[1], parsedCommand[2])
		response = "OK\n"
	case "rpop":
		if len(parsedCommand) != 2 {
			response = "ERR wrong number of arguments for 'rpop' command\nRequired `rpop <key>`\n"
			err = errors.New("ERR wrong number of arguments for 'rpop' command\nRequired `rpop <key>`")
			break
		}
		value, ok := data.RPop(parsedCommand[1])
		if ok {
			response = "\"" + value + "\"" + "\n"
		} else {
			response = "(nil)\n"
		}
	case "llen":
		if len(parsedCommand) != 2 {
			response = "ERR wrong number of arguments for 'llen' command\nRequired `llen <key>`\n"
			err = errors.New("ERR wrong number of arguments for 'llen' command\nRequired `llen <key>`")
			break
		}
		length := data.LLen(parsedCommand[1])
		response = strconv.Itoa(length) + "\n"
	case "exit":
		response = "Bye\n"
		err = errors.New("exit")
	default:
		response = "ERR unknown command '" + cmd + "'"
	}
	return response, err
}
