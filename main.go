package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bogevold/task-mgr/internal/handler"
	pgstore "github.com/bogevold/task-mgr/internal/store"
	"github.com/bogevold/task-mgr/internal/task"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	store, err := pgstore.NewPostgresStore(connStr)
	if err != nil {
		fmt.Printf("Kunne ikke initialisere postgres backend: %v", err)
		return
	}
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
