package main

import (
	"path/filepath"
	"testing"
)

func TestAuthServiceRegisterAndLogin(t *testing.T) {
	origUsersPath := usersFilePath
	usersFilePath = filepath.Join(t.TempDir(), "users.json")
	t.Cleanup(func() {
		usersFilePath = origUsersPath
	})

	authService := NewAuthService()

	// 1. Yeni kullanıcı kaydı
	token, user, err := authService.Register("testuser", "secret123")
	if err != nil {
		t.Fatalf("Register başarısız: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("Beklenen kullanıcı adı testuser, gelen: %s", user.Username)
	}
	if token == "" {
		t.Errorf("Geçerli token dönmeliydi")
	}

	// 2. Token doğrulama
	authenticatedUser, ok := authService.Authenticate(token)
	if !ok {
		t.Fatalf("Oluşturulan token doğrulanamadı")
	}
	if authenticatedUser.ID != user.ID {
		t.Errorf("Doğrulanan kullanıcı ID uyuşmuyor")
	}

	// 3. Tekrar aynı kullanıcı adıyla kayıt olmaya çalışma
	_, _, err = authService.Register("testuser", "secret123")
	if err == nil {
		t.Errorf("Aynı kullanıcı adı ile ikinci kez kayıt olunamamalı")
	}

	// 4. Doğru şifre ile giriş
	loginToken, loginUser, err := authService.Login("testuser", "secret123")
	if err != nil {
		t.Fatalf("Doğru şifreyle login başarısız: %v", err)
	}
	if loginUser.ID != user.ID || loginToken == "" {
		t.Errorf("Login sonucu beklenmedik")
	}

	// 5. Yanlış şifre ile giriş
	_, _, err = authService.Login("testuser", "wrongpassword")
	if err == nil {
		t.Errorf("Yanlış şifre ile login hata dönmeliydi")
	}
}

