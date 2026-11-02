package main

import (
	"errors"
	"strings"
	"sync"
	"time"
)

type AuthService struct {
	tokensMu sync.RWMutex
	tokens   map[string]UserSession
}

type UserSession struct {
	User      User
	ExpiresAt time.Time
}

func NewAuthService() *AuthService {
	return &AuthService{
		tokens: make(map[string]UserSession),
	}
}

func (s *AuthService) Login(username, password string) (string, User, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", User{}, errors.New("kullanıcı adı ve şifre zorunludur")
	}

	user, err := findUserByUsername(username)
	if err != nil {
		return "", User{}, errors.New("kullanıcı bilgileri okunamadı")
	}
	if user == nil || !user.CheckPassword(password) {
		return "", User{}, errors.New("geçersiz kullanıcı adı veya şifre")
	}

	token := generateRandomHex(24)
	s.tokensMu.Lock()
	s.tokens[token] = UserSession{
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.tokensMu.Unlock()

	return token, *user, nil
}

func (s *AuthService) Register(username, password string) (string, User, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", User{}, errors.New("kullanıcı adı ve şifre zorunludur")
	}

	existing, err := findUserByUsername(username)
	if err != nil {
		return "", User{}, errors.New("kullanıcı kontrolü yapılamadı")
	}
	if existing != nil {
		return "", User{}, errors.New("bu kullanıcı adı zaten kullanılıyor")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	newUser := User{
		ID:        generateRandomHex(16),
		Username:  username,
		Role:      RoleUser,
		CreatedAt: now,
	}
	newUser.SetPassword(password)

	if err := newUser.Validate(); err != nil {
		return "", User{}, err
	}

	if err := saveNewUser(newUser); err != nil {
		return "", User{}, errors.New("kullanıcı kaydedilemedi")
	}

	token := generateRandomHex(24)
	s.tokensMu.Lock()
	s.tokens[token] = UserSession{
		User:      newUser,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.tokensMu.Unlock()

	return token, newUser, nil
}

func (s *AuthService) Authenticate(token string) (User, bool) {
	if token == "" {
		return User{}, false
	}

	s.tokensMu.RLock()
	session, exists := s.tokens[token]
	s.tokensMu.RUnlock()

	if !exists {
		return User{}, false
	}

	if time.Now().After(session.ExpiresAt) {
		s.tokensMu.Lock()
		delete(s.tokens, token)
		s.tokensMu.Unlock()
		return User{}, false
	}

	return session.User, true
}

