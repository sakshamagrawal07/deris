package data

var Data map[string]string

func InitData() {
	Data = make(map[string]string)
}

func Set(key string, value string) {
	Data[key] = value
}

func Get(key string) (string, bool) {
	value, ok := Data[key]
	return value, ok
}
