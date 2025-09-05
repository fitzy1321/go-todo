package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	db    *db.AppDB
	state sessionState

	entryForm tea.Model
	table     tea.Model
	extra     string

	keyMap AppKeyMap
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
	m.extra = ""
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
					return m, nil
				}

				id, err := mTable.SelectedId()
				if err != nil {
					m.extra = err.Error()
					return m, nil
				}
				if id == nil {
					m.extra = "Warn: 'delete cmd' could not find todo object"
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
					return m, nil
				}

				id, err := mTable.SelectedId()
				if err != nil {
					m.extra = err.Error()
					return m, nil
				}
				if id == nil {
					m.extra = "Warn: 'toggle cmd' could not find todo object"
					return m, nil
				}
				m.table, cmd = mTable.Update(ToggleMsg{id: *id})
				return m, cmd
			}
		case msg.String() == "enter":
			if m.state == entryFormView {
				// Get value from form
				mTable := m.table.(TodoTableModel)
				mEntryForm := m.entryForm.(EntryFormModel)
				title := mEntryForm.Value()
				if title == "" {
					return m, nil
				}
				// Make a new Todo
				if err := mTable.AddTodo(title); err != nil {
					// TODO: do something here
				}
				mEntryForm = NewEntryForm()
				mEntryForm.Blur()
				m.entryForm = mEntryForm
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

func (m AppModel) View() string {
	switch m.state {
	case tableView:
		res := "Golang Todo TUI\n\n" + m.table.View()
		if m.extra != "" {
			res += "\n" + m.extra + "\n"
		}
		return res

	case entryFormView:
		return m.entryForm.View()
	default:
		return "Something has gone wrong ..."
	}
}
