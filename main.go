package main

import (
	"log"

	"todo.tui.go/app"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	path, err := app.TodosFilePath()
	if err != nil {
		log.Fatalf("failed to prepare todo storage: %v", err)
	}

	todos, err := app.LoadTodos(path)
	if err != nil {
		log.Fatalf("failed to load todos: %v", err)
	}

	p := tea.NewProgram(app.InitialModel(path, todos), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("app exited with error: %v", err)
	}
}
