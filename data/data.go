package data

var Data map[string][]string
var expire int

func InitData() {
	Data = make(map[string][]string)
	expire = 0
}

func Set(key string, value string) {
	Data[key] = []string{value}
}

func Setnx(key string, value string) {
	if _, ok := Data[key]; !ok {
		Data[key] = []string{value}
	}
}

func Get(key string) ([]string, bool) {
	value, ok := Data[key]
	return value, ok
}

func LPush(key string,value string){
	Data[key] = append([]string{value},Data[key]...)
}

func LPop(key string) (string, bool) {
	if len(Data[key]) == 0 {
		return "", false
	}
	value := Data[key][0]
	Data[key] = Data[key][1:]
	return value, true
}

func RPush(key string,value string){
	Data[key] = append(Data[key], value)
}

func RPop(key string) (string, bool) {
	if len(Data[key]) == 0 {
		return "", false
	}
	value := Data[key][len(Data[key])-1]
	Data[key] = Data[key][:len(Data[key])-1]
	return value, true
}

func LLen(key string) int {
	return len(Data[key])
}