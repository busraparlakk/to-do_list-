package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var todosMu sync.Mutex

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	request.Title = strings.TrimSpace(request.Title)

	if request.Title == "" {
		http.Error(w, "title boş olamaz", http.StatusBadRequest)
		return
	}

	if request.DueDate != "" {
		_, err := time.Parse(time.RFC3339, request.DueDate)
		if err != nil {
			http.Error(w, "due_date geçerli bir RFC3339 tarihi olmalı", http.StatusBadRequest)
			return
		}
	}

	idBytes := make([]byte, 16)
	_, err = rand.Read(idBytes)
	if err != nil {
		log.Println("ID üretilemedi:", err)
		http.Error(w, "Todo oluşturulamadı", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)

	todo := Todo{
		ID:          hex.EncodeToString(idBytes),
		Title:       request.Title,
		Description: request.Description,
		Done:        false,
		DueDate:     request.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	todosMu.Lock()
	defer todosMu.Unlock()

	todos, err := readTodos()
	if err != nil {
		log.Println("Todo dosyası okunamadı:", err)
		http.Error(w, "Todo oluşturulamadı", http.StatusInternalServerError)
		return
	}

	todos = append(todos, todo)

	err = writeTodos(todos)
	if err != nil {
		log.Println("Todo dosyasına yazılamadı:", err)
		http.Error(w, "Todo kaydedilemedi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	todosMu.Lock()
	todos, err := readTodos()
	todosMu.Unlock()

	if err != nil {
		log.Println("Todo dosyası okunamadı:", err)
		http.Error(w, "Todolar alınamadı", http.StatusInternalServerError)
		return
	}

	if todos == nil {
		todos = []Todo{}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}
func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	todosMu.Lock()
	todos, err := readTodos()
	todosMu.Unlock()

	if err != nil {
		log.Println("Todo dosyası okunamadı:", err)
		http.Error(w, "Todo alınamadı", http.StatusInternalServerError)
		return
	}

	for _, todo := range todos {
		if todo.ID == id {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(todo)
			return
		}
	}

	http.Error(w, "Todo bulunamadı", http.StatusNotFound)
}