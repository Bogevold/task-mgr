package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bogevold/task-mgr/internal/handler"
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

	th := handler.NewTaskHandler(store)
	mux := http.NewServeMux()
	th.RegisterRoutes(mux)
	http.ListenAndServe(":8072", mux)
}
