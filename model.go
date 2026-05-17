package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type appMode string

const (
	modeList appMode = "list"
	modeAdd  appMode = "add"
)

type model struct {
	todos    []Todo
	cursor   int
	input    textinput.Model
	mode     appMode
	filePath string
	err      error
}

func initialModel(filePath string, todos []Todo) model {
	input := textinput.New()
	input.Placeholder = "Add a todo and press Enter"
	input.CharLimit = 200
	input.Width = 50

	return model{
		todos:    todos,
		cursor:   0,
		input:    input,
		mode:     modeList,
		filePath: filePath,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.mode == modeAdd {
			switch msg.String() {
			case "enter":
				title := strings.TrimSpace(m.input.Value())
				if title != "" {
					m.todos = append(m.todos, Todo{Title: title, Done: false})
					m.cursor = len(m.todos) - 1
					m = m.persist()
				}
				m.mode = modeList
				m.input.SetValue("")
				m.input.Blur()
				return m, nil
			case "esc":
				m.mode = modeList
				m.input.SetValue("")
				m.input.Blur()
				return m, nil
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case "n":
			m.mode = modeAdd
			m.input.SetValue("")
			m.input.Focus()
		case " ":
			if len(m.todos) > 0 {
				m.todos[m.cursor].Done = !m.todos[m.cursor].Done
				m = m.persist()
			}
		case "x":
			if len(m.todos) > 0 {
				m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
				if m.cursor >= len(m.todos) && m.cursor > 0 {
					m.cursor--
				}
				m = m.persist()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	if m.mode == modeAdd {
		b.WriteString("New todo: ")
		b.WriteString(m.input.View())
		b.WriteString("\n\n")
	}

	if len(m.todos) == 0 {
		b.WriteString("No todos yet. Press 'n' to add one.\n")
	} else {
		b.WriteString("Todos:\n")
		for i, todo := range m.todos {
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}

			status := "[ ]"
			if todo.Done {
				status = "[x]"
			}

			b.WriteString(fmt.Sprintf("%s %s %s\n", cursor, status, todo.Title))
		}
	}

	b.WriteString("\nShortcuts: j/k or arrows move | n new | space toggle done | x archive(delete) | q quit\n")

	if m.err != nil {
		b.WriteString("Error: ")
		b.WriteString(m.err.Error())
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) persist() model {
	m.err = SaveTodos(m.filePath, m.todos)
	return m
}
