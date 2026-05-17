package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	appDirName = "todo-tui-go"
	todosFile  = "todos.json"
)

func todosFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(configDir, appDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, todosFile)
	if err := ensureFile(path); err != nil {
		return "", err
	}

	return path, nil
}

func ensureFile(path string) error {
	_, err := os.Stat(path)

	// if there is no error: file exists, do nothing
	if err == nil {
		return nil
	}
	// if the error is not "file does not exist", return error
	if !os.IsNotExist(err) {
		return err
	}
	// if the file does not exist, create it with an empty array
	return os.WriteFile(path, []byte("[]\n"), 0o644)
}

func LoadTodos(path string) ([]Todo, error) {
	data, err := os.ReadFile(path)

	// if there is an error reading the file, return error
	if err != nil {
		return nil, err
	}

	// if the file is empty, return an empty slice
	if len(data) == 0 {
		return []Todo{}, nil
	}

	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}

	if todos == nil {
		return []Todo{}, nil
	}

	return todos, nil
}

func SaveTodos(path string, todos []Todo) error {
	if todos == nil {
		todos = []Todo{}
	}

	data, err := json.MarshalIndent(todos, "", "  ")

	// if there is an error, return error
	if err != nil {
		return err
	}

	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
