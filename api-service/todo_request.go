package main

type CreateTodoRequest struct { // todo varken başka bir yapı oluşturduk çünkü
	Title       string `json:"title"`       // todo kaydın tamamını temsil eder
	Description string `json:"description"` // CreateTodoRequest ise kullanıcının oluşturma isteğinde göndereceği alanları temsiö eder.
	DueDate     string `json:"due_date"`
}
