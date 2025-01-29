package commands

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sakshamagrawal07/deris/config"
	"github.com/sakshamagrawal07/deris/data"
	"github.com/sakshamagrawal07/deris/utils"
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

func RecoverFromAOF() {
	cmds, err := utils.ReadStringDataFromFile(config.LogFile)
	if err != nil {
		log.Println(err)
		return
	}
	for _, cmd := range cmds {
		QueueCommand(cmd, -1)
	}
}

func ClearAof() {
	filename := config.LogFile

	err := os.Remove(filename)
	if err != nil {
		log.Println("Error clearing aof ", err)
		return
	}
	log.Println("AOF Cleared")
}

func QueueCommand(cmd string, fd int) {
	log.Println("Queueing normally ", cmd)
	data.CommandQueue.Push(cmd, fd)
}

func MultiCommandQueue(cmd string, fd int) {
	log.Println("Queueing normally multi ", cmd)
	command := parseCommand(cmd)
	if len(command) == 1 && (command[0] == DISCARD_CMD || command[0] == EXEC_CMD) {
		resp, err := ExecuteCommand(cmd, fd)
		if err != nil {
			log.Println("Error executing command ", err)
			utils.RespondToClientWithFd(fd, err.Error())
		} else {
			utils.RespondToClientWithFd(fd, resp)
		}
		return
	}
	data.MultiCommandQueue.Push(cmd, fd)
}

func ExecuteCommandsInQueue() {
	// cmds := data.CommandQueue
	log.Println("Command Queue Execution started")
	for {
		if !data.CommandQueue.IsEmpty() {
			cmd, fd := data.CommandQueue.Pop()
			resp, err := ExecuteCommand(cmd, fd)
			if err != nil {
				log.Println("Error executing command: ", err)
			}
			// unix.Write(clientFd, []byte(response))
			if fd > 0 {
				log.Println("Writing to ", fd)
				utils.RespondToClientWithFd(fd, resp)
			}
			if resp == "Bye\n" {
				// Client closed connection
				utils.CloseClientConnection(fd)
			}
		}
	}
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

	inputCommand := cmd

	cmd = parsedCommand[0]

	if config.AppendOnly && utils.Contains(WriteCommands, cmd) {
		utils.AppendDataToFileAsString(config.LogFile, inputCommand)
	}

	if config.MultiCommand && cmd != DISCARD_CMD && cmd != EXEC_CMD {
		return "queued", nil
	}

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
	case MULTI_CMD:
		if len(parsedCommand) != 1 {
			setErrorMessage(&response, &err, MULTI_CMD, MULTI_CMD_FORMAT)
			break
		}
		config.MultiCommand = true
		response = "OK\n"
	case DISCARD_CMD:
		if len(parsedCommand) != 1 {
			setErrorMessage(&response, &err, DISCARD_CMD, DISCARD_CMD_FORMAT)
			break
		}
		config.MultiCommand = false
		data.MultiCommandQueue.Clear()
		response = "OK\n"
	case EXEC_CMD:
		if len(parsedCommand) != 1 {
			setErrorMessage(&response, &err, EXEC_CMD, EXEC_CMD_FORMAT)
			break
		}
		config.MultiCommand = false
		data.CommandQueue.Copy(&data.MultiCommandQueue)
		response = "OK\n"
	case ASYNC_SAVE_DB:
		if len(parsedCommand) != 1 {
			setErrorMessage(&response, &err, ASYNC_SAVE_DB, ASYNC_SAVE_DB_FORMAT)
			break
		}
		go data.AsyncSaveDB()
		response = "OK\n"
	case EXIT:
		response = "Bye\n"

	default:
		response = "ERR unknown command '" + cmd + "'"
	}

	return response, err
}
