package data

import (
	"errors"
	"time"

	"github.com/sakshamagrawal07/deris/config"
	"github.com/sakshamagrawal07/deris/utils"
)

var Data map[string]utils.DataStruct
var expiringKeysTree utils.RadixNode
// var CmdQueue bool
var CommandQueue utils.Queue
var MultiCommandQueue utils.Queue

func InitData() {
	Data = make(map[string]utils.DataStruct)
	expiringKeysTree = *utils.NewRadixNode("", time.Time{}, false)
	CommandQueue.Init()
	MultiCommandQueue.Init()
	// CmdQueue = false
}

func Get(key string) ([]string, bool) {
	if expireTime, found := expiringKeysTree.Find(key); found {
		if expireTime.Before(time.Now()) || expireTime == time.Now() {
			Delete(key)
			return nil, false
		}
	}
	if structBody, ok := Data[key]; ok {
		return structBody.Values, ok
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

func DeleteExpiredNodes() {
	expiringKeysTree.DeleteExpiredNodes()
}

func BackupData() {
	utils.WriteDataToFileAsJson(config.DataFilePath,Data)
}

func Delete(key string) {
	delete(Data, key)
	expiringKeysTree.Delete(key)
}

func Set(key string, value string) {
	Data[key] = *utils.NewDataStruct([]string{value}, false, []int{})
}

func Setnx(key string, value string) {
	if _, ok := Data[key]; !ok {
		Data[key] = *utils.NewDataStruct([]string{value}, false, []int{})
	}
}

func LPush(key string, value string) {
	structBody, ok := Data[key]
	if ok {
		structBody.Values = append([]string{value},structBody.Values...)
		Data[key] = structBody
		return
	}
	Data[key] = *utils.NewDataStruct(append(Data[key].Values, value), false, []int{})
}

func LPop(key string) (string, bool) {
	if len(Data[key].Values) == 0 {
		return "", false
	}

	structBody := Data[key]

	value := structBody.Values[0]
	structBody.Values = structBody.Values[1:]
	Data[key] = structBody
	return value, true
}

func RPush(key string, value string) {
	structBody,ok := Data[key]
	if ok {
		structBody.Values = append(structBody.Values, value)
		Data[key] = structBody
		return
	}
	Data[key] = *utils.NewDataStruct(append(Data[key].Values,value),false,[]int{})
}

func RPop(key string) (string, bool) {
	structBody := Data[key]
	if len(structBody.Values) == 0 {
		return "", false
	}
	len := len(structBody.Values) - 1
	value := structBody.Values[len]
	structBody.Values = structBody.Values[:len]
	Data[key] = structBody
	return value, true
}

func LLen(key string) int {
	return len(Data[key].Values)
}

func SubscribeToKey(key string, fd int) {
	if structBody, ok := Data[key]; ok {
		structBody.Subscribed = true
		structBody.Subscribers = append(structBody.Subscribers, fd)
		Data[key] = structBody
		return
	}
	Data[key] = *utils.NewDataStruct([]string{}, true, []int{fd})
}

func UnsubscribeToKey(key string, fd int) {
	if structBody, ok := Data[key]; ok {
		structBody.Subscribers = utils.RemoveFromIntSlice(structBody.Subscribers, fd)
		if len(structBody.Subscribers) == 0 {
			structBody.Subscribed = false
		}
		Data[key] = structBody
	}
}

func PublishToKey(key string, value string) error {
	if structBody, ok := Data[key]; ok {
		RPush(key, value)
		for _, subscriber := range structBody.Subscribers {
			utils.RespondToClientWithFd(subscriber, value)
		}
		return nil
	}
	return errors.New("key ('" + key + "') not found")
}
