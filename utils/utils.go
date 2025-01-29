package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

var Clients map[int]struct{}

func CloseClientConnection(fd int) {
	if fd <= 0 {
		return
	}
	unix.Close(fd)
	delete(Clients, fd)
	fmt.Printf("Connection closed: %d\n", fd)
}

func RemoveFromIntSlice(slice []int, val int) []int {
	for x, v := range slice {
		if v == val {
			return append(slice[:x], slice[x+1:]...)
		}
	}
	return slice
}

func RespondToClientWithFd(clientFd int, response string) error {
	if _, err := unix.Write(clientFd, []byte(response)); err != nil {
		return err
	}
	return nil
}

func WriteDataToFileAsJson(filename string, data map[string]DataStruct) error {
	log.Println("Backup data cron : ", time.Now())

	if data == nil {
		log.Println("Attempted to write data in : ",filename)
		return errors.New("map is empty")
	}
	// dir,_ := os.Getwd()
	// log.Println("Current Working Directory : ", dir)
	// filename := config.DataFilePath
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	// log.Println("File open")

	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	// log.Println("Data converted")

	_, err = file.Write(dataJSON)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Data successfully written to %s\n.", filename)
	return nil
}

func AppendDataToFileAsString(filename string, val string) error {
	// Open the file in append mode, create if not exists
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.WriteString(val) // Adding a newline for better readability
	if err != nil {
		log.Println("Error writing to file:", err)
		return err
	}

	log.Printf("Data successfully appended to %s\n", filename)
	return nil
}

func ReadJsonDataFromFile(filename string, data map[string]DataStruct) error {
	// filename := config.DataFilePath
	// Check if the file exists
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return errors.New("file does not exist")
	}

	// Read the file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	log.Println("Data successfully read from file.")
	return nil
}

func ReadStringDataFromFile(filename string) ([]string, error) {
	// Check if the file exists
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("file does not exist")
	}

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read lines into a slice of strings
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func Contains(list []string, str string) bool {
	for _, val := range list {
		if val == str {
			return true
		}
	}
	return false
}
