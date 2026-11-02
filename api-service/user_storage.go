package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	usersFilePath = "data/users.json"
	usersMu       sync.Mutex
)

func initUserStorage() error {
	usersMu.Lock()
	defer usersMu.Unlock()

	if _, err := os.Stat(usersFilePath); errors.Is(err, os.ErrNotExist) {
		return seedDefaultUsers()
	}

	data, err := os.ReadFile(usersFilePath)
	if err != nil {
		return err
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil || len(users) == 0 {
		return seedDefaultUsers()
	}

	return nil
}

func seedDefaultUsers() error {
	now := time.Now().UTC().Format(time.RFC3339)

	admin := User{
		ID:        generateRandomHex(16),
		Username:  "admin",
		Role:      RoleSuperAdmin,
		CreatedAt: now,
	}
	admin.SetPassword("admin123")

	user1 := User{
		ID:        generateRandomHex(16),
		Username:  "user1",
		Role:      RoleUser,
		CreatedAt: now,
	}
	user1.SetPassword("user123")

	return saveUsersUnsafe([]User{admin, user1})
}

func readUsers() ([]User, error) {
	usersMu.Lock()
	defer usersMu.Unlock()

	data, err := os.ReadFile(usersFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := seedDefaultUsers(); err != nil {
				return nil, err
			}
			data, err = os.ReadFile(usersFilePath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func saveUsersUnsafe(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(usersFilePath, data, 0644)
}

func findUserByUsername(username string) (*User, error) {
	users, err := readUsers()
	if err != nil {
		return nil, err
	}

	normalized := strings.ToLower(strings.TrimSpace(username))
	for _, u := range users {
		if strings.ToLower(u.Username) == normalized {
			return &u, nil
		}
	}

	return nil, nil
}

func saveNewUser(user User) error {
	usersMu.Lock()
	defer usersMu.Unlock()

	data, err := os.ReadFile(usersFilePath)
	var users []User
	if err == nil {
		_ = json.Unmarshal(data, &users)
	}

	users = append(users, user)
	return saveUsersUnsafe(users)
}

func generateRandomHex(byteCount int) string {
	b := make([]byte, byteCount)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

