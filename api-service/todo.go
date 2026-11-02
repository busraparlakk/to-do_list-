package main

import (
	"errors"
	"strings"
	"time"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

type Todo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Done        bool       `json:"done"`
	DueDate     string     `json:"due_date"`
	Visibility  string     `json:"visibility"`
	UserID      string     `json:"user_id"`
	Username    string     `json:"username,omitempty"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// Normalize, kullanıcıdan gelen alanları düzenler ve varsayılanları atar.
func (todo *Todo) Normalize() {
	todo.Title = strings.TrimSpace(todo.Title)
	todo.Visibility = strings.ToLower(strings.TrimSpace(todo.Visibility))
	if todo.Visibility == "" {
		todo.Visibility = VisibilityPrivate
	}
}

// Validate, todo'nun geçerliliğini ve iş kurallarını denetler.
func (todo Todo) Validate() error {
	if todo.Title == "" {
		return errors.New("title boş olamaz")
	}

	if todo.Visibility != VisibilityPublic && todo.Visibility != VisibilityPrivate {
		return errors.New("visibility 'public' veya 'private' olmalı")
	}

	if todo.DueDate != "" {
		if _, err := time.Parse(time.RFC3339, todo.DueDate); err != nil {
			return errors.New("due_date geçerli bir RFC3339 tarihi olmalı")
		}
	}

	return nil
}

// CanView, kullanıcının todo'yu görüntüleme yetkisi olup olmadığını kontrol eder.
func (todo Todo) CanView(user User) bool {
	if user.IsSuperAdmin() {
		return true
	}
	if todo.IsPublic() || todo.Visibility == "" {
		return true
	}
	return todo.UserID == user.ID
}

// CanModify, kullanıcının todo'yu güncelleme veya silme yetkisini kontrol eder.
func (todo Todo) CanModify(user User) bool {
	if user.IsSuperAdmin() {
		return true
	}
	if todo.UserID == "" {
		return true
	}
	return todo.UserID == user.ID
}

// IsPublic, todo'nun herkese açık olup olmadığını belirtir.
func (todo Todo) IsPublic() bool {
	return todo.Visibility == VisibilityPublic
}

// IsDeleted, kaydın silinmiş olup olmadığını belirtir.
func (todo Todo) IsDeleted() bool {
	return todo.DeletedAt != nil
}

// MarkDeleted, kaydı soft delete yöntemiyle silinmiş olarak işaretler.
func (todo *Todo) MarkDeleted() {
	if todo.IsDeleted() {
		return
	}

	now := time.Now().UTC()
	todo.DeletedAt = &now
	todo.UpdatedAt = now.Format(time.RFC3339)
}
