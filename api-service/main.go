package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	if err := initUserStorage(); err != nil {
		log.Println("Kullanıcı veri tabanı başlatılırken uyarı:", err)
	}

	authService := NewAuthService()
	authHandler := NewAuthHandler(authService)

	mux := http.NewServeMux()

	// Genel endpoint'ler
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)

	// Kimlik doğrulaması gerektiren endpoint'ler
	withAuth := func(h http.HandlerFunc) http.HandlerFunc {
		return authMiddleware(authService, h)
	}

	mux.HandleFunc("GET /api/auth/me", withAuth(authHandler.Me))
	mux.HandleFunc("POST /api/todos", withAuth(createTodoHandler))
	mux.HandleFunc("GET /api/todos", withAuth(listTodosHandler))
	mux.HandleFunc("GET /api/todos/{id}", withAuth(getTodoHandler))
	mux.HandleFunc("PUT /api/todos/{id}", withAuth(updateTodoHandler))
	mux.HandleFunc("DELETE /api/todos/{id}", withAuth(deleteTodoHandler))

	log.Println("API servisi http://localhost:4000 adresinde çalışıyor")

	err := http.ListenAndServe(":4000", corsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
