package utils

import "golang.org/x/sys/unix"

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
