package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task-api/middleware"
	"task-api/model"
	"task-api/store"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	Store *store.MemoryStore
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	var req model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	b := make([]byte, 8)
	rand.Read(b)
	
	task := &model.Task{
		ID:          fmt.Sprintf("%x", b),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	h.Store.CreateTask(task)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if limit <= 0 { limit = 10 }
	if page <= 0 { page = 1 }
	offset := (page - 1) * limit

	tasks, total := h.Store.ListTasks(userID, limit, offset)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  tasks,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	task, err := h.Store.GetTask(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	var req model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	task, err := h.Store.GetTask(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status

	h.Store.UpdateTask(task)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	if err := h.Store.DeleteTask(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
