package model

type TaskState int64

const (
	Open TaskState = iota
	Cancelled
	Done
)

var mapTaskState = map[TaskState]string{
	Open:      " ",
	Cancelled: "-",
	Done:      "x",
}

func (s TaskState) Trigger() string {
	return mapTaskState[s]
}
