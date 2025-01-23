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

func setErrorMessage(response* string,err* error,command string,commandFormat string){
	*response = "ERR wrong number of arguments for '" + command + "' command\nRequired `" + commandFormat + "`\n"
	*err = errors.New("ERR wrong number of arguments for '" + command + "' command\nRequired `" + commandFormat)
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
	case GET:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response,&err,GET,GET_FORMAT)
			break
		}
		value, ok := data.Get(parsedCommand[1])
		if ok {
			response = strings.Join(value, " ") + "\n"
		} else {
			response = "(nil)\n"
		}

	case STRING_SET:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response,&err,STRING_SET,STRING_SET_FORMAT)
			break
		}
		data.Set(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case STRING_SET_NOT_EXISTS:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response,&err,STRING_SET_NOT_EXISTS,STRING_SET_NOT_EXISTS_FORMAT)
			break
		}
		data.Setnx(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case LIST_LEFT_PUSH:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response,&err,LIST_LEFT_PUSH,LIST_LEFT_PUSH_FORMAT)
			break
		}
		data.LPush(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case LIST_LEFT_POP:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response,&err,LIST_LEFT_POP,LIST_LEFT_POP_FORMAT)
			break
		}
		value, ok := data.LPop(parsedCommand[1])
		if ok {
			response = "\"" + value + "\"" + "\n"
		} else {
			response = "(nil)\n"
		}

	case LIST_RIGHT_PUSH:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response,&err,LIST_RIGHT_PUSH,LIST_RIGHT_PUSH_FORMAT)
			break
		}
		data.RPush(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case LIST_RIGHT_POP:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response,&err,LIST_RIGHT_POP,LIST_RIGHT_POP_FORMAT)
			break
		}
		value, ok := data.RPop(parsedCommand[1])
		if ok {
			response = "\"" + value + "\"" + "\n"
		} else {
			response = "(nil)\n"
		}

	case LIST_LENGTH:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response,&err,LIST_LENGTH,LIST_LENGTH_FORMAT)
			break
		}
		length := data.LLen(parsedCommand[1])
		response = strconv.Itoa(length) + "\n"

	case EXIT:
		response = "Bye\n"
		err = errors.New("exit")

	default:
		response = "ERR unknown command '" + cmd + "'"
	}

	return response, err
}