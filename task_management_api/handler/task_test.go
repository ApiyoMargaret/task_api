package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"task-api/handler"
	"task-api/middleware"
	"task-api/model"
	"task-api/store"
)

func TestCreateTask_Validation(t *testing.T) {
	db := store.NewMemoryStore()
	h := &handler.TaskHandler{Store: db}

	// Title is too short, status invalid
	body, _ := json.Marshal(model.TaskRequest{Title: "ab", Status: "invalid-status"})
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	
	// Inject Mock Context Auth
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user123")
	req = req.WithContext(ctx)
	
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422, got %d", rec.Code)
	}
}
