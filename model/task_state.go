package model

type TaskState int64

const (
	Open TaskState = iota
	Cancelled
	Done
)

var taskStateNames = [...]string{"open", "cancelled", "done"}

func (s TaskState) String() string {
	if int(s) >= len(taskStateNames) {
		return ""
	}
	return taskStateNames[s]
}

func (s *TaskState) MarshalJSON() ([]byte, error) {
	return []byte("\"" + s.String() + "\""), nil
}

var mapTaskState = map[TaskState]string{
	Open:      " ",
	Cancelled: "-",
	Done:      "x",
}

func (s TaskState) Trigger() string {
	return mapTaskState[s]
}
