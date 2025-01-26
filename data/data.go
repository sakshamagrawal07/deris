package data

import (
	"errors"
	"time"

	"github.com/sakshamagrawal07/deris/utils"
)

var Data map[string]dataStruct
var expiringKeysTree utils.RadixNode

func InitData() {
	Data = make(map[string]dataStruct)
	expiringKeysTree = *utils.NewRadixNode("", time.Time{}, false)
}

func Get(key string) ([]string, bool) {
	if expireTime, found := expiringKeysTree.Find(key); found {
		if expireTime.Before(time.Now()) || expireTime == time.Now() {
			Delete(key)
			return nil, false
		}
	}
	if structBody, ok := Data[key]; ok {
		return structBody.values, ok
	}
	return nil, false
}

func Expire(key string, seconds int) {
	if seconds < 0 {
		Delete(key)
		return
	}
	expirationTime := time.Now().Add(time.Second * time.Duration(seconds))
	expiringKeysTree.Insert(key, expirationTime)
}

func deleteExpiredNodes() {
	expiringKeysTree.DeleteExpiredNodes()
}

func Delete(key string) {
	delete(Data, key)
	expiringKeysTree.Delete(key)
}

func Set(key string, value string) {
	Data[key] = *newDataStruct([]string{value}, false, []int{})
}

func Setnx(key string, value string) {
	if _, ok := Data[key]; !ok {
		Data[key] = *newDataStruct([]string{value}, false, []int{})
	}
}

func LPush(key string, value string) {
	if structBody, ok := Data[key]; ok {
		structBody.values = append(structBody.values, value)
	}
	Data[key] = *newDataStruct(append(Data[key].values, value), false, []int{})
}

func LPop(key string) (string, bool) {
	if len(Data[key].values) == 0 {
		return "", false
	}

	structBody := Data[key]

	value := structBody.values[0]
	structBody.values = structBody.values[1:]
	Data[key] = structBody
	return value, true
}

func RPush(key string, value string) {
	structBody := Data[key]
	structBody.values = append(structBody.values, value)
	Data[key] = structBody
}

func RPop(key string) (string, bool) {
	structBody := Data[key]
	if len(structBody.values) == 0 {
		return "", false
	}
	len := len(structBody.values) - 1
	value := structBody.values[len]
	structBody.values = structBody.values[:len]
	Data[key] = structBody
	return value, true
}

func LLen(key string) int {
	return len(Data[key].values)
}

func SubscribeToKey(key string, fd int) {
	if structBody, ok := Data[key]; ok {
		structBody.subscribed = true
		structBody.subscribers = append(structBody.subscribers, fd)
		Data[key] = structBody
		return
	}
	Data[key] = *newDataStruct([]string{}, true, []int{fd})
}

func UnsubscribeToKey(key string, fd int) {
	if structBody, ok := Data[key]; ok {
		structBody.subscribers = utils.RemoveFromIntSlice(structBody.subscribers, fd)
		if len(structBody.subscribers) == 0 {
			structBody.subscribed = false
		}
		Data[key] = structBody
	}
}

func PublishToKey(key string, value string) error {
	if structBody, ok := Data[key]; ok {
		RPush(key, value)
		for _, subscriber := range structBody.subscribers {
			utils.RespondToClientWithFd(subscriber, value)
		}
		return nil
	}
	return errors.New("key ('" + key + "') not found")
}
