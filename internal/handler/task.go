package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bogevold/task-mgr/internal/task"
)

type TaskHandler struct {
	store task.TaskStore
}

type UpdateTaskRequest struct {
	ID        *uint      `json:"id"`
	Title     *string    `json:"title"`
	Done      *bool      `json:"done"`
	CreatedAt *time.Time `json:"created_at"`
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

func (h *TaskHandler) handleGetById(w http.ResponseWriter, r *http.Request) {
	str_id := r.PathValue("id")

	id, err := strconv.ParseUint(str_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := h.store.Get(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handleSave(w http.ResponseWriter, r *http.Request) {
	// les task fra r, lagre, skriv svar til w
	var newTask task.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newTask.CreatedAt = time.Now()
	savedTask, err := h.store.Save(newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedTask)
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	str_id := r.PathValue("id")

	id, err := strconv.ParseUint(str_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.store.Delete(uint(id))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	str_id := r.PathValue("id")
	id, err := strconv.ParseUint(str_id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task2Update, err := h.store.Get(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var reqTask UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&reqTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqTask.Title != nil {
		task2Update.Title = *reqTask.Title
	}

	if reqTask.Done != nil {
		task2Update.Done = *reqTask.Done
	}

	updatedTask, err := h.store.Update(task2Update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.handleGetAll)
	mux.HandleFunc("GET /tasks/{id}", h.handleGetById)
	mux.HandleFunc("POST /tasks", h.handleSave)
	mux.HandleFunc("DELETE /tasks/{id}", h.handleDelete)
	mux.HandleFunc("PATCH /tasks/{id}", h.handleUpdate)
}
