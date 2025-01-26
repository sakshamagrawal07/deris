package data

type dataStruct struct {
	values      []string
	subscribed  bool
	subscribers []int
}

func newDataStruct(values []string, subscribed bool, fds []int) *dataStruct {
	return &dataStruct{
		values:      values,
		subscribed:  subscribed,
		subscribers: fds,
	}
}
