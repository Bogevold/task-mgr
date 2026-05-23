package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bogevold/task-mgr/internal/task"
)

type TaskHandler struct {
	store task.TaskStore
}

func NewTaskHandler(store task.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// func (h *TaskHandler) handleSave(w http.ResponseWriter, r *http.Request) {
// 	// les task fra r, lagre, skriv svar til w
// 	newTask
// 	json.NewDecoder(r.Body).Decode(newTask)
// 	savedTask, err := h.store.Save(newTask)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	json.NewEncoder(w).Encode(savedTask)
// }

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleGetAll)
}
