package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"task-api/middleware"
	"task-api/model"
	"task-api/store"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Store *store.MemoryStore
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	
	b := make([]byte, 8)
	rand.Read(b)
	id := fmt.Sprintf("%x", b)

	user := &model.User{ID: id, Email: req.Email, Password: string(hashed)}
	if err := h.Store.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	user, err := h.Store.GetUserByEmail(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := middleware.GenerateToken(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
