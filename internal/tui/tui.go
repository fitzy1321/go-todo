package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fitzy1321/go-todo/internal/db"
	"github.com/google/uuid"
)

type sessionState int

const (
	tableView sessionState = iota
	entryFormView
)

type AppKeyMap struct {
	TodoTableKeyMap
}

type AppModel struct {
	db       *db.AppDB
	errorStr string
	keyMap   AppKeyMap
	state    sessionState

	// Nested Models
	entryForm tea.Model
	table     tea.Model
}

func NewApp(db *db.AppDB) *AppModel {
	addTodoForm := NewEntryForm()
	table := NewTodoTable(db)
	table.Focus()
	t := AppModel{
		db:        db,
		state:     tableView,
		entryForm: addTodoForm,
		table:     table,
		keyMap:    AppKeyMap{DefaultTodoTableKeyMap()},
	}
	return &t
}

type (
	TableMsg     struct{}
	EntryFormMsg struct{}
	ToggleMsg    struct{ id uuid.UUID }
	DeleteMsg    struct{ id uuid.UUID }
)

func (t *ToggleMsg) Id() uuid.UUID {
	return t.id
}

func (t *DeleteMsg) Id() uuid.UUID {
	return t.id
}

func (m AppModel) Init() tea.Cmd { return nil }

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.errorStr = ""
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case EntryFormMsg:
		m.state = entryFormView
	case TableMsg:
		m.state = tableView

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Help):
			if m.state == tableView {
				m.table, cmd = m.table.Update(msg)
				return m, cmd
			}
		case key.Matches(msg, m.keyMap.New):
			// Change state to Show EntryForm
			if m.state == tableView {
				mTable := m.table.(TodoTableModel)
				mTable.Blur()
				m.table = mTable

				mEntryForm := m.entryForm.(EntryFormModel)
				mEntryForm.Focus()
				m.entryForm = mEntryForm

				return m, func() tea.Msg { return EntryFormMsg{} }
			}
		case key.Matches(msg, m.keyMap.Delete):
			// Delete Todo
			if m.state == tableView {
				mTable := m.table.(TodoTableModel)
				if len(mTable.todos) == 0 {
					m.errorStr = "'delete': No items to delete"
					return m, nil
				}

				id := mTable.SelectedId()
				if id == nil {
					m.errorStr = "'delete': could not find selected row's id"
					return m, nil
				}
				m.table, cmd = mTable.Update(DeleteMsg{id: *id})
				return m, cmd
			}
		case key.Matches(msg, m.keyMap.Toggle):
			// Toggle Todo.Completed
			if m.state == tableView {
				mTable := m.table.(TodoTableModel)
				if len(mTable.todos) == 0 {
					m.errorStr = "'toggle': No items to toggle"
					return m, nil
				}

				id := mTable.SelectedId()
				if id == nil {
					m.errorStr = "'toggle': could not find selected row's"
					return m, nil
				}
				m.table, cmd = mTable.Update(ToggleMsg{id: *id})
				return m, cmd
			}
		case msg.String() == "enter":
			if m.state == entryFormView {
				// Get value from form
				mEntryForm := m.entryForm.(EntryFormModel)
				newTitle := mEntryForm.Value()
				if newTitle == "" {
					return m, nil
				}
				mEntryForm = NewEntryForm()
				mEntryForm.Blur()
				m.entryForm = mEntryForm

				// Make a new Todo
				mTable := m.table.(TodoTableModel)
				mTable.AddTodo(newTitle)
				mTable.Focus()
				m.table = mTable

				// send tea.Cmd to update
				return m, func() tea.Msg { return TableMsg{} }
			}
		}
	}
	switch m.state {
	case tableView:
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	case entryFormView:
		m.entryForm, cmd = m.entryForm.Update(msg)
		return m, cmd
	}
	return m, nil
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B50235"))
)

func (m AppModel) View() string {
	switch m.state {
	case tableView:
		var b strings.Builder
		b.WriteString(titleStyle.Render("Golang Todo TUI") + "\n\n" + m.table.View() + "\n")
		if m.errorStr != "" {
			b.WriteString(errorStyle.Render(m.errorStr))
		}
		return b.String()

	case entryFormView:
		return m.entryForm.View()
	default:
		return "Something has gone wrong ..."
	}
}
