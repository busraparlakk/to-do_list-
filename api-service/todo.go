package main

import (
	"errors"
	"strings"
	"time"
)

type Todo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Done        bool       `json:"done"`
	DueDate     string     `json:"due_date"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// Normalize, kullanıcıdan gelen metni düzenler.
func (todo *Todo) Normalize() {
	todo.Title = strings.TrimSpace(todo.Title)
}

// Validate, todo'nun ortak iş kurallarına uygunluğunu kontrol eder.
func (todo Todo) Validate() error {
	if todo.Title == "" {
		return errors.New("title boş olamaz")
	}

	if todo.DueDate != "" {
		if _, err := time.Parse(time.RFC3339, todo.DueDate); err != nil {
			return errors.New("due_date geçerli bir RFC3339 tarihi olmalı")
		}
	}

	return nil
}

// IsDeleted, kaydın silinmiş olarak işaretlenip işaretlenmediğini bildirir.
func (todo Todo) IsDeleted() bool {
	return todo.DeletedAt != nil
}

// MarkDeleted, kaydı silinmiş olarak işaretler.
// Zaten silinmişse ilk silinme tarihini korur.
func (todo *Todo) MarkDeleted() {
	if todo.IsDeleted() {
		return
	}

	now := time.Now().UTC()

	todo.DeletedAt = &now
	todo.UpdatedAt = now.Format(time.RFC3339)
}
