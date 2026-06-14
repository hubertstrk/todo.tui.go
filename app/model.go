package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appMode string

const (
	modeList appMode = "list"
	modeAdd  appMode = "add"
	modeEdit appMode = "edit"
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

var doneTodoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4dffa3"))
var appFrameStyle = lipgloss.NewStyle().
	// BorderStyle(lipgloss.RoundedBorder()).
	// BorderForeground(lipgloss.Color("240")).
	Padding(0, 1)

var promptLabelStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#a94dff")).
	Bold(true)

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("204"))

func (x todoItem) Title() string {
	status := "☐"
	if x.todo.Done {
		status = "✓"
		return doneTodoStyle.Render(status + " " + x.todo.Title)
	}

	return status + " " + x.todo.Title
}

func (x todoItem) Description() string {
	return ""
}

func (x todoItem) FilterValue() string {
	return x.todo.Title
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
	input.Width = 200

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	todoList := list.New(createItems(todos), delegate, 0, 0)
	todoList.Title = "Todos"
	todoList.SetShowStatusBar(true)
	todoList.SetFilteringEnabled(false)
	todoList.SetShowHelp(true)

	todoList.Help.Styles.ShortKey = todoList.Help.Styles.ShortKey.
		Foreground(lipgloss.Color("#d9d9d9")).
		Padding(0, 1)

	todoList.Help.Styles.FullKey = todoList.Help.Styles.FullKey.
		Foreground(lipgloss.Color("#d9d9d9")).
		Padding(0, 1)

	newKey := key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new todo"))
	toggleKey := key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle"))
	archiveKey := key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete"))
	editKey := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))

	todoList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{newKey, toggleKey, archiveKey, editKey}
	}
	todoList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{newKey, toggleKey, archiveKey, editKey}
	}

	return model{
		todos:    todos,
		list:     todoList,
		input:    input,
		mode:     modeList,
		filePath: filePath,
	}
}

// Init initializes the model and returns an optional command to execute.
// Called when the program starts, before the first update and view calls.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages (events) and updates the model accordingly.
// Called whenever an event occurs (e.g., key press, window resize).
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		listHeight := msg.Height - 6
		if listHeight < 1 {
			listHeight = 1
		}
		listWidth := msg.Width - 4
		if listWidth < 1 {
			listWidth = 1
		}
		m.list.SetSize(listWidth, listHeight)
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

		if m.mode == modeEdit {
			switch msg.String() {
			case "enter":
				title := strings.TrimSpace(m.input.Value())
				if title != "" {
					idx := m.list.Index()
					m.todos[idx].Title = title
					m.list.SetItems(createItems(m.todos))
					m.list.Select(idx)
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
		case "e":
			if len(m.todos) > 0 {
				idx := m.list.Index()
				m.input.SetValue(m.todos[idx].Title)
				m.mode = modeEdit
				m.input.Focus()
				return m, nil
			}
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

// View renders the UI based on the current model state.
// Called after every update to refresh the display.
func (m model) View() string {
	var b strings.Builder

	switch m.mode {
	case modeAdd:
		b.WriteString(promptLabelStyle.Render("New todo"))
		b.WriteString("\n")
		b.WriteString(m.input.View())
	case modeEdit:
		b.WriteString(promptLabelStyle.Render("Edit todo"))
		b.WriteString("\n")
		b.WriteString(m.input.View())
	default:
		b.WriteString(m.list.View())
	}

	if m.err != nil {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	return appFrameStyle.Render(b.String())
}

func (m model) persist() model {
	m.err = SaveTodos(m.filePath, m.todos)
	return m
}
