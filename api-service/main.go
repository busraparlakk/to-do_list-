package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/todos", createTodoHandler)
	mux.HandleFunc("GET /api/todos", listTodosHandler)
	mux.HandleFunc("GET /api/todos/{id}", getTodoHandler)
	mux.HandleFunc("PUT /api/todos/{id}", updateTodoHandler)
	mux.HandleFunc("DELETE /api/todos/{id}", deleteTodoHandler)

	log.Println("API servisi http://localhost:4000 adresinde çalışıyor")

	err := http.ListenAndServe(":4000", corsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
