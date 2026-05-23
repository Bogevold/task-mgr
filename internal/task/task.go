package task

import (
	"time"
)

type Task struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}
