package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
	list     list.Model
	input    textinput.Model
	mode     appMode
	filePath string
	err      error
}

type todoItem struct {
	todo Todo
}

func (i todoItem) Title() string {
	status := "[ ]"
	if i.todo.Done {
		status = "[x]"
	}

	return status + " " + i.todo.Title
}

func (i todoItem) Description() string {
	return ""
}

func (i todoItem) FilterValue() string {
	return i.todo.Title
}

func createItems(todos []Todo) []list.Item {
	items := make([]list.Item, 0, len(todos))
	for _, todo := range todos {
		items = append(items, todoItem{todo: todo})
	}

	return items
}

func InitialModel(filePath string, todos []Todo) tea.Model {
	input := textinput.New()
	input.Placeholder = "Add a todo and press Enter"
	input.CharLimit = 200
	input.Width = 50

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	todoList := list.New(createItems(todos), delegate, 0, 0)
	todoList.Title = "Todos"
	todoList.SetShowStatusBar(true)
	todoList.SetFilteringEnabled(false)
	todoList.SetShowHelp(true)

	newKey := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new todo"))
	toggleKey := key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle"))
	archiveKey := key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete"))
	todoList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{newKey, toggleKey, archiveKey}
	}
	todoList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{newKey, toggleKey, archiveKey}
	}

	return model{
		todos:    todos,
		list:     todoList,
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
	case tea.WindowSizeMsg:
		listHeight := msg.Height - 4
		if listHeight < 1 {
			listHeight = 1
		}
		m.list.SetSize(msg.Width, listHeight)
		return m, nil
	case tea.KeyMsg:
		if m.mode == modeAdd {
			switch msg.String() {
			case "enter":
				title := strings.TrimSpace(m.input.Value())
				if title != "" {
					m.todos = append(m.todos, Todo{Title: title, Done: false})
					m.list.SetItems(createItems(m.todos))
					m.list.Select(len(m.todos) - 1)
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
		case "n":
			m.mode = modeAdd
			m.input.SetValue("")
			m.input.Focus()
			return m, nil
		case " ":
			if len(m.todos) > 0 {
				idx := m.list.Index()
				m.todos[idx].Done = !m.todos[idx].Done
				m.list.SetItems(createItems(m.todos))
				m.list.Select(idx)
				m = m.persist()
			}
			return m, nil
		case "x":
			if len(m.todos) > 0 {
				idx := m.list.Index()
				m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
				m.list.SetItems(createItems(m.todos))
				if len(m.todos) == 0 {
					m.list.Select(0)
				} else if idx >= len(m.todos) {
					m.list.Select(len(m.todos) - 1)
				} else {
					m.list.Select(idx)
				}
				m = m.persist()
			}
			return m, nil
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	if m.mode == modeAdd {
		b.WriteString("New todo: ")
		b.WriteString(m.input.View())
		b.WriteString("\n\n")
	} else {
		b.WriteString(m.list.View())
	}

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
