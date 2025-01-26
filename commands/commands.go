package commands

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/sakshamagrawal07/deris/data"
)

// func ReadCommand(c net.Conn) (string, error) {
// 	var buf []byte = make([]byte, 512)
// 	n, err := c.Read(buf[:])
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(buf[:n]), nil
// }

func parseCommand(cmd string) []string {
	log.Println("Parsing command...")
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	return strings.Fields(cmd)
}

func setErrorMessage(response *string, err *error, command string, commandFormat string) {
	*response = "ERR wrong number of arguments for '" + command + "' command\nRequired `" + commandFormat + "`\n"
	*err = errors.New("ERR wrong number of arguments for '" + command + "' command\nRequired `" + commandFormat)
}

func ExecuteCommand(cmd string, fd int) (string, error) {
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
			setErrorMessage(&response, &err, GET, GET_FORMAT)
			break
		}
		value, ok := data.Get(parsedCommand[1])
		if ok {
			response = strings.Join(value, " ") + "\n"
		} else {
			response = "(nil)\n"
		}

	case DELETE_KEY:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response, &err, DELETE_KEY, DELETE_KEY_FORMAT)
			break
		}
		data.Delete(parsedCommand[1])
		response = "OK\n"

	case EXPIRE_KEY:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response, &err, EXPIRE_KEY, EXPIRE_KEY_FORMAT)
			break
		}
		seconds, err := strconv.Atoi(parsedCommand[2])
		if err != nil {
			setErrorMessage(&response, &err, EXPIRE_KEY, EXPIRE_KEY_FORMAT)
			break
		}
		data.Expire(parsedCommand[1], seconds)
		response = "OK\n"

	case STRING_SET:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response, &err, STRING_SET, STRING_SET_FORMAT)
			break
		}
		data.Set(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case STRING_SET_NOT_EXISTS:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response, &err, STRING_SET_NOT_EXISTS, STRING_SET_NOT_EXISTS_FORMAT)
			break
		}
		data.Setnx(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case LIST_LEFT_PUSH:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response, &err, LIST_LEFT_PUSH, LIST_LEFT_PUSH_FORMAT)
			break
		}
		data.LPush(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case LIST_LEFT_POP:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response, &err, LIST_LEFT_POP, LIST_LEFT_POP_FORMAT)
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
			setErrorMessage(&response, &err, LIST_RIGHT_PUSH, LIST_RIGHT_PUSH_FORMAT)
			break
		}
		data.RPush(parsedCommand[1], parsedCommand[2])
		response = "OK\n"

	case LIST_RIGHT_POP:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response, &err, LIST_RIGHT_POP, LIST_RIGHT_POP_FORMAT)
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
			setErrorMessage(&response, &err, LIST_LENGTH, LIST_LENGTH_FORMAT)
			break
		}
		length := data.LLen(parsedCommand[1])
		response = strconv.Itoa(length) + "\n"

	case SUBSCRIBE_KEY:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response, &err, SUBSCRIBE_KEY, SUBSCRIBE_KEY_FORMAT)
			break
		}
		data.SubscribeToKey(parsedCommand[1], fd)
		response = "OK\n"

	case UNSUBSCRIBE_KEY:
		if len(parsedCommand) != 2 {
			setErrorMessage(&response, &err, UNSUBSCRIBE_KEY, UNSUBSCRIBE_KEY_FORMAT)
			break
		}
		data.UnsubscribeToKey(parsedCommand[1], fd)
		response = "OK\n"

	case PUBLISH_KEY:
		if len(parsedCommand) != 3 {
			setErrorMessage(&response, &err, PUBLISH_KEY, PUBLISH_KEY_FORMAT)
			break
		}
		err = data.PublishToKey(parsedCommand[1], parsedCommand[2])
		if err != nil {
			response = "(nil)\n"
		} else {
			response = "OK\n"
		}
	case EXIT:
		response = "Bye\n"
		
	default:
		response = "ERR unknown command '" + cmd + "'"
	}

	return response, err
}
