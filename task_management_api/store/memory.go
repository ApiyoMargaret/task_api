package store

import (
	"errors"
	"sync"
	"task-api/model"
)

type MemoryStore struct {
	mu     sync.RWMutex
	users  map[string]*model.User
	tasks  map[string]*model.Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users: make(map[string]*model.User),
		tasks: make(map[string]*model.Task),
	}
}

func (s *MemoryStore) CreateUser(u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[u.Email]; exists {
		return errors.New("user already exists")
	}
	s.users[u.Email] = u
	return nil
}

func (s *MemoryStore) GetUserByEmail(email string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, exists := s.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (s *MemoryStore) CreateTask(t *model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
}

func (s *MemoryStore) GetTask(id, userID string) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, exists := s.tasks[id]
	if !exists || t.UserID != userID {
		return nil, errors.New("task not found")
	}
	return t, nil
}

func (s *MemoryStore) ListTasks(userID string, limit, offset int) ([]*model.Task, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var userTasks []*model.Task
	for _, t := range s.tasks {
		if t.UserID == userID {
			userTasks = append(userTasks, t)
		}
	}

	total := len(userTasks)
	if offset > total {
		return []*model.Task{}, total
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return userTasks[offset:end], total
}

func (s *MemoryStore) UpdateTask(t *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[t.ID]; !exists {
		return errors.New("task not found")
	}
	s.tasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTask(id, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, exists := s.tasks[id]
	if !exists || t.UserID != userID {
		return errors.New("task not found")
	}
	delete(s.tasks, id)
	return nil
}