package task

import (
	"time"
)

type Task struct {
	ID        uint
	Title     string
	Done      bool
	CreatedAt time.Time
}
