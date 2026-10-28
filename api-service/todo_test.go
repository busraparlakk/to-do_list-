package main

import "testing"

func TestTodoNormalizeAndValidate(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		dueDate   string
		wantTitle string
		wantError string
	}{
		{
			name:      "geçerli başlık",
			title:     "Go öğren",
			wantTitle: "Go öğren",
		},
		{
			name:      "başlığın kenarındaki boşluklar temizlenir",
			title:     "  Go öğren  ",
			wantTitle: "Go öğren",
		},
		{
			name:      "boş başlık reddedilir",
			title:     "",
			wantTitle: "",
			wantError: "title boş olamaz",
		},
		{
			name:      "yalnızca boşluk içeren başlık reddedilir",
			title:     " \t\n ",
			wantTitle: "",
			wantError: "title boş olamaz",
		},
		{
			name:      "geçerli tarih kabul edilir",
			title:     "Go öğren",
			dueDate:   "2026-11-01T18:00:00+03:00",
			wantTitle: "Go öğren",
		},
		{
			name:      "tarih verilmemesi kabul edilir",
			title:     "Go öğren",
			dueDate:   "",
			wantTitle: "Go öğren",
		},
		{
			name:      "geçersiz tarih reddedilir",
			title:     "Go öğren",
			dueDate:   "yarın",
			wantTitle: "Go öğren",
			wantError: "due_date geçerli bir RFC3339 tarihi olmalı",
		},
		{
			name:      "olmayan takvim tarihi reddedilir",
			title:     "Go öğren",
			dueDate:   "2026-02-30T18:00:00Z",
			wantTitle: "Go öğren",
			wantError: "due_date geçerli bir RFC3339 tarihi olmalı",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := Todo{
				Title:   tt.title,
				DueDate: tt.dueDate,
			}

			todo.Normalize()

			if todo.Title != tt.wantTitle {
				t.Fatalf(
					"Başlık: beklenen %q, gelen %q",
					tt.wantTitle,
					todo.Title,
				)
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
				t.Fatalf(
					"Hata mesajı: beklenen %q, gelen %q",
					tt.wantError,
					err.Error(),
				)
			}
		})
	}
}
