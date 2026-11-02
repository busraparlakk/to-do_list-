package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	RoleUser       = "user"
	RoleSuperAdmin = "superadmin"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	CreatedAt    string `json:"created_at"`
}

func (u User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

func (u *User) SetPassword(plainPassword string) {
	u.PasswordHash = hashPassword(plainPassword)
}

func (u User) CheckPassword(plainPassword string) bool {
	return u.PasswordHash == hashPassword(plainPassword)
}

func (u User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("kullanıcı adı boş olamaz")
	}
	if len(u.Username) < 3 {
		return errors.New("kullanıcı adı en az 3 karakter olmalı")
	}
	if u.PasswordHash == "" {
		return errors.New("şifre boş olamaz")
	}
	if u.Role != RoleUser && u.Role != RoleSuperAdmin {
		return errors.New("geçersiz rol")
	}
	return nil
}

func hashPassword(password string) string {
	salt := "todo_salt_key_2026"
	hash := sha256.Sum256([]byte(salt + password))
	return hex.EncodeToString(hash[:])
}

