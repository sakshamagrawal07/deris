package utils

import (
	"fmt"
)

type queueStruct struct {
	cmd string
	fd  int
}

type Queue struct {
	data []queueStruct
}

func (q *Queue) Init() {
	q.data = []queueStruct{}
}

func (q *Queue) Push(val string, fd int) {
	q.data = append(q.data, queueStruct{cmd: val, fd: fd})
	// log.Println("Current Queue After Push : ")
	// q.Print()
}

func (q *Queue) Pop() (string, int) {
	if len(q.data) == 0 {
		return "empty", 0
	}
	val := q.data[0]
	q.data = q.data[1:]
	// log.Println("Current Queue After Push : ")
	// q.Print()
	return val.cmd, val.fd
}

func (q *Queue) IsEmpty() bool {
	return q.Length() <= 0
}

func (q *Queue) Print() {
	var found bool = false
	for x, val := range q.data {
		found = true
		fmt.Println(x, val)
	}
	if !found {
		fmt.Println("Queue empty.")
	}
}

func (q *Queue) Length() int {
	return len(q.data)
}

func (q *Queue) Clear() {
	q.data = []queueStruct{}
}

func (q1 *Queue) Copy(q2 *Queue) {
	for !q2.IsEmpty() {
		q1.Push(q2.Pop())
	}
}
