package main

import "testing"

func TestTodoNormalizeAndValidate(t *testing.T) {
	tests := []struct {
		name           string
		title          string
		dueDate        string
		visibility     string
		wantTitle      string
		wantVisibility string
		wantError      string
	}{
		{
			name:           "geçerli başlık ve boş visibility varsayılan private olur",
			title:          "Go öğren",
			visibility:     "",
			wantTitle:      "Go öğren",
			wantVisibility: VisibilityPrivate,
		},
		{
			name:           "public visibility kabul edilir",
			title:          "Public görev",
			visibility:     "public",
			wantTitle:      "Public görev",
			wantVisibility: VisibilityPublic,
		},
		{
			name:           "geçersiz visibility reddedilir",
			title:          "Görev",
			visibility:     "invalid_vis",
			wantTitle:      "Görev",
			wantVisibility: "invalid_vis",
			wantError:      "visibility 'public' veya 'private' olmalı",
		},
		{
			name:           "başlığın kenarındaki boşluklar temizlenir",
			title:          "  Go öğren  ",
			wantTitle:      "Go öğren",
			wantVisibility: VisibilityPrivate,
		},
		{
			name:           "boş başlık reddedilir",
			title:          "",
			wantTitle:      "",
			wantVisibility: VisibilityPrivate,
			wantError:      "title boş olamaz",
		},
		{
			name:           "geçerli tarih kabul edilir",
			title:          "Go öğren",
			dueDate:        "2026-11-01T18:00:00+03:00",
			wantTitle:      "Go öğren",
			wantVisibility: VisibilityPrivate,
		},
		{
			name:           "geçersiz tarih reddedilir",
			title:          "Go öğren",
			dueDate:        "yarın",
			wantTitle:      "Go öğren",
			wantVisibility: VisibilityPrivate,
			wantError:      "due_date geçerli bir RFC3339 tarihi olmalı",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := Todo{
				Title:      tt.title,
				DueDate:    tt.dueDate,
				Visibility: tt.visibility,
			}

			todo.Normalize()

			if todo.Title != tt.wantTitle {
				t.Fatalf("Başlık: beklenen %q, gelen %q", tt.wantTitle, todo.Title)
			}
			if todo.Visibility != tt.wantVisibility {
				t.Fatalf("Visibility: beklenen %q, gelen %q", tt.wantVisibility, todo.Visibility)
			}

			err := todo.Validate()
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("Hata beklenmiyordu: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Hata bekleniyordu: %s", tt.wantError)
			}

			if err.Error() != tt.wantError {
				t.Fatalf("Hata mesajı: beklenen %q, gelen %q", tt.wantError, err.Error())
			}
		})
	}
}

func TestTodoAuthorizationRules(t *testing.T) {
	superAdmin := User{ID: "admin-id", Role: RoleSuperAdmin, Username: "admin"}
	user1 := User{ID: "user1-id", Role: RoleUser, Username: "user1"}
	user2 := User{ID: "user2-id", Role: RoleUser, Username: "user2"}

	privateTodo := Todo{
		ID:         "t1",
		Title:      "User1 Gizli Not",
		Visibility: VisibilityPrivate,
		UserID:     user1.ID,
	}

	publicTodo := Todo{
		ID:         "t2",
		Title:      "User1 Açık Not",
		Visibility: VisibilityPublic,
		UserID:     user1.ID,
	}

	// 1. Private Todo Görüntüleme Kuralları
	if !privateTodo.CanView(user1) {
		t.Errorf("Sahip (user1) kendi private todo'sunu görebilmeli")
	}
	if privateTodo.CanView(user2) {
		t.Errorf("Başka kullanıcı (user2) user1'in private todo'sunu görememeli")
	}
	if !privateTodo.CanView(superAdmin) {
		t.Errorf("Süper admin tüm private todo'ları görebilmeli")
	}

	// 2. Public Todo Görüntüleme Kuralları
	if !publicTodo.CanView(user1) {
		t.Errorf("Sahip (user1) kendi public todo'sunu görebilmeli")
	}
	if !publicTodo.CanView(user2) {
		t.Errorf("Başka kullanıcı (user2) public todo'yu görebilmeli")
	}
	if !publicTodo.CanView(superAdmin) {
		t.Errorf("Süper admin public todo'yu görebilmeli")
	}

	// 3. Güncelleme ve Silme (CanModify) Kuralları
	if !privateTodo.CanModify(user1) {
		t.Errorf("Sahip kendi todo'sunu güncelleyebilmeli")
	}
	if privateTodo.CanModify(user2) {
		t.Errorf("Başka kullanıcı başkasının todo'sunu güncelleyememeli")
	}
	if !publicTodo.CanModify(user1) {
		t.Errorf("Sahip public todo'sunu güncelleyebilmeli")
	}
	if publicTodo.CanModify(user2) {
		t.Errorf("Başka kullanıcı public olsa dahi başkasının todo'sunu güncelleyememeli")
	}
	if !privateTodo.CanModify(superAdmin) {
		t.Errorf("Süper admin tüm todo'ları yönetebilmeli")
	}
}
