package utils

type DataStruct struct {
	Values      []string
	Subscribed  bool
	Subscribers []int
}

func NewDataStruct(values []string, subscribed bool, fds []int) *DataStruct {
	return &DataStruct{
		Values:      values,
		Subscribed:  subscribed,
		Subscribers: fds,
	}
}
