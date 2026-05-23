package main

import (
	"fmt"
	"time"

	"github.com/bogevold/task-mgr/internal/task"
)

func main() {
	store := task.NewInMemoryStore()
	title := ""
	for i := range 4 {
		title = fmt.Sprintf("Task %d", i)
		var tsk = task.Task{Title: title, CreatedAt: time.Now()}
		_, err := store.Save(tsk)
		if err != nil {
			fmt.Printf("forventet ingen feil, fikk: %v", err)
		}
	}

	tasks, err := store.GetAll()
	if err != nil {
		fmt.Printf("Kunne ikke hente oppgaver: %v", err)
	}
	for _, task := range tasks {
		fmt.Printf("ID: %d, Title: %s, Done: %v\n", task.ID, task.Title, task.Done)
	}

}
