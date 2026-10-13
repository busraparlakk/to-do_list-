package main

import (
	"path/filepath"
	"testing"
)

func TestWriteAndReadTodos(t *testing.T) {
	originalPath := todosFilePath
	todosFilePath = filepath.Join(t.TempDir(), "todos.json")

	t.Cleanup(func() {
		todosFilePath = originalPath
	})

	todos := []Todo{
		{
			ID:    "1",
			Title: "Go ile dosya işlemlerini öğren",
			Done:  false,
		},
	}

	err := writeTodos(todos)
	if err != nil {
		t.Fatalf("Dosyaya yazılamadı: %v", err)
	}

	savedTodos, err := readTodos()
	if err != nil {
		t.Fatalf("Dosya okunamadı: %v", err)
	}

	if len(savedTodos) != 1 {
		t.Fatalf("1 kayıt bekleniyordu, gelen: %d", len(savedTodos))
	}

	if savedTodos[0] != todos[0] {
		t.Fatalf("Okunan kayıt, yazılan kayıtla aynı değil")
	}

	t.Log("Todo dosyaya yazıldı ve aynı içerikle geri okundu")
}
