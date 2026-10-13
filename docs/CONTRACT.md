# Veri modeli — bir todo nesnesinin alanları ve tipleri
id: string              benzersiz kimlik                    
title: string           zorunlu boş olamaz
description: string     opsiyonel
done: boolean           varsayılan: false
due_date: string        ISO 8601 formatında, opsiyonel
created_at: string      oluşturulma zamanı
updated_at: string      güncelleme zamanı

# endpointler
    
GET /todos
Request Body: Yok.
Response Body:
[
  {
    "id": "todo-123",
    "title": "Reverse proxy öğren",
    "description": "Go ile reverse proxy geliştir",
    "done": false,
    "due_date": "2026-09-30T18:00:00Z",
    "created_at": "2026-09-25T15:00:00Z",
    "updated_at": "2026-09-25T15:00:00Z"
  }
]

Status Codes: 200 OK, 500 Internal Server Error

GET /todos/{id}
Request Body: Yok.
Response Body:
{
  "id": "todo-123",
  "title": "Reverse proxy öğren",
  "description": "Go ile reverse proxy geliştir",
  "done": false,
  "due_date": "2026-09-30T18:00:00Z",
  "created_at": "2026-09-25T15:00:00Z",
  "updated_at": "2026-09-25T15:00:00Z"
}

Status Codes: 200 OK, 404 Not Found, 500 Internal Server Error

POST /todos
Request Body:
{
  "title": "Reverse proxy öğren",
  "description": "Go ile reverse proxy geliştir",
  "due_date": "2026-09-30T18:00:00Z"
}

Response Body:
{
  "id": "todo-123",
  "title": "Reverse proxy öğren",
  "description": "Go ile reverse proxy geliştir",
  "done": false,
  "due_date": "2026-09-30T18:00:00Z",
  "created_at": "2026-09-25T15:00:00Z",
  "updated_at": "2026-09-25T15:00:00Z"
}

Status Codes: 201 Created, 400 Bad Request, 500 Internal Server Error

PUT /todos/{id}
Request Body:
{
  "title": "Reverse proxy çalış",
  "description": "Go ile reverse proxy geliştir",
  "done": true,
  "due_date": "2026-09-30T18:00:00Z"
}

Response Body:
{
  "id": "todo-123",
  "title": "Reverse proxy çalış",
  "description": "Go ile reverse proxy geliştir",
  "done": true,
  "due_date": "2026-09-30T18:00:00Z",
  "created_at": "2026-09-25T15:00:00Z",
  "updated_at": "2026-09-25T16:00:00Z"
}

Status Codes: 200 OK, 400 Bad Request, 404 Not Found, 500 Internal Server Error

DELETE /todos/{id}
Request Body: Yok.
Response Body: Yok.

Status Codes: 204 No Content, 404 Not Found, 500 Internal Server Error