package model

import (
	"fmt"
)

type Summary struct {
	Done  int `json:"done"`
	Total int `json:"total"`
}

func (s Summary) String() string {
	return fmt.Sprintf("%d/%d tasks", s.Done, s.Total)
}
