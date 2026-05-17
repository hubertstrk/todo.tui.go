package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	path, err := todosFilePath()
	if err != nil {
		log.Fatalf("failed to prepare todo storage: %v", err)
	}

	todos, err := LoadTodos(path)
	if err != nil {
		log.Fatalf("failed to load todos: %v", err)
	}

	p := tea.NewProgram(initialModel(path, todos))
	if _, err := p.Run(); err != nil {
		log.Fatalf("app exited with error: %v", err)
	}
}
