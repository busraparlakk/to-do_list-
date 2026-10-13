package main

import (
	"encoding/json"
	"os"
)

var todosFilePath = "data/todos.json"

func readTodos() ([]Todo, error) {
	data, err := os.ReadFile(todosFilePath)
	if err != nil {
		return nil, err
	}

	var todos []Todo

	err = json.Unmarshal(data, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func writeTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(todosFilePath, data, 0644)
}
