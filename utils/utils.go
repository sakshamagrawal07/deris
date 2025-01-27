package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

const DataFilePath = "./diskStorage/data.txt"

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

func WriteDataToFile(data map[string]DataStruct) error {
	log.Println("Backup data cron : ", time.Now())
	// dir,_ := os.Getwd()
	// log.Println("Current Working Directory : ", dir)
	filename := DataFilePath
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	log.Println("File open")

	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Data converted")

	_, err = file.Write(dataJSON)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Data successfully written to file.")
	return nil
}

func ReadDataFromFile(data map[string]DataStruct) error {
	filename := DataFilePath
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
