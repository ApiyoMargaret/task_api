package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Task struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // e.g., "pending", "completed"
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (t TaskRequest) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Title, validation.Required, validation.Length(3, 100)),
		validation.Field(&t.Status, validation.In("pending", "completed")),
	)
}