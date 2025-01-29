package config

const DataFilePath = "./data.txt"

const LogFile = "./aof.txt"

var Host string
var Port int
var ExpireKeyCronTimer int
var BackupCronTimer int

var ClearAOF bool

var MultiCommand bool = false

var AppendOnly bool = false