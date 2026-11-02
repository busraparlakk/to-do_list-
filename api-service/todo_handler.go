package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var todosMu sync.Mutex

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
		return
	}

	var request CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
		return
	}

	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
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
		Visibility:  request.Visibility,
		UserID:      user.ID,
		Username:    user.Username,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	todo.Normalize()
	if err := todo.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
	if err := writeTodos(todos); err != nil {
		log.Println("Todo dosyasına yazılamadı:", err)
		http.Error(w, "Todo kaydedilemedi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(todo)
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
		return
	}

	todosMu.Lock()
	todos, err := readTodos()
	todosMu.Unlock()

	if err != nil {
		log.Println("Todo dosyası okunamadı:", err)
		http.Error(w, "Todolar alınamadı", http.StatusInternalServerError)
		return
	}

	activeTodos := make([]Todo, 0)
	for _, todo := range todos {
		if !todo.IsDeleted() && todo.CanView(user) {
			activeTodos = append(activeTodos, todo)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(activeTodos)
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
		return
	}

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
		if todo.ID == id && !todo.IsDeleted() {
			if !todo.CanView(user) {
				http.Error(w, "Bu todoyu görüntüleme yetkiniz yok", http.StatusForbidden)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(todo)
			return
		}
	}

	http.Error(w, "Todo bulunamadı", http.StatusNotFound)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	var request UpdateTodoRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Geçersiz JSON veya alan", http.StatusBadRequest)
		return
	}

	if err := decoder.Decode(new(any)); err != io.EOF {
		http.Error(w, "Body tek bir JSON nesnesi içermeli", http.StatusBadRequest)
		return
	}

	if !hasUpdateFields(request) {
		http.Error(w, "En az bir güncellenecek alan gönderilmeli", http.StatusBadRequest)
		return
	}

	todosMu.Lock()
	defer todosMu.Unlock()

	todos, err := readTodos()
	if err != nil {
		log.Println("Todo dosyası okunamadı:", err)
		http.Error(w, "Todo güncellenemedi", http.StatusInternalServerError)
		return
	}

	for i := range todos {
		if todos[i].ID != id || todos[i].IsDeleted() {
			continue
		}

		if !todos[i].CanModify(user) {
			http.Error(w, "Bu todoyu güncelleme yetkiniz yok", http.StatusForbidden)
			return
		}

		applyTodoUpdates(&todos[i], request)
		todos[i].Normalize()

		if err := todos[i].Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todos[i].UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		if err := writeTodos(todos); err != nil {
			log.Println("Todo dosyasına yazılamadı:", err)
			http.Error(w, "Todo kaydedilemedi", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(todos[i])
		return
	}

	http.Error(w, "Todo bulunamadı", http.StatusNotFound)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Yetkisiz erişim", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	todosMu.Lock()
	defer todosMu.Unlock()

	todos, err := readTodos()
	if err != nil {
		log.Println("Todo dosyası okunamadı:", err)
		http.Error(w, "Todo silinemedi", http.StatusInternalServerError)
		return
	}

	for i := range todos {
		if todos[i].ID != id || todos[i].IsDeleted() {
			continue
		}

		if !todos[i].CanModify(user) {
			http.Error(w, "Bu todoyu silme yetkiniz yok", http.StatusForbidden)
			return
		}

		todos[i].MarkDeleted()
		if err := writeTodos(todos); err != nil {
			log.Println("Todo dosyasına yazılamadı:", err)
			http.Error(w, "Todo silinemedi", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Error(w, "Todo bulunamadı", http.StatusNotFound)
}

func hasUpdateFields(req UpdateTodoRequest) bool {
	return req.Title != nil || req.Description != nil ||
		req.Done != nil || req.DueDate != nil || req.Visibility != nil
}

func applyTodoUpdates(todo *Todo, req UpdateTodoRequest) {
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Done != nil {
		todo.Done = *req.Done
	}
	if req.DueDate != nil {
		todo.DueDate = *req.DueDate
	}
	if req.Visibility != nil {
		todo.Visibility = *req.Visibility
	}
}
