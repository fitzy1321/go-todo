package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fitzy1321/go-todo/internal/db"
	"github.com/google/uuid"
)

type sessionState int

const (
	tableView sessionState = iota
	entryFormView
)

type AppModel struct {
	db    *db.AppDB
	state sessionState

	entryForm tea.Model
	table     tea.Model
	extra     string
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
	}
	return &t
}

type (
	EntryFormMsg struct{}
	TableMsg     struct{}
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "n":
			// Make a new Todo
			if m.state == tableView {
				mTable := m.table.(TodoTableModel)
				mTable.Blur()
				m.table = mTable
				mEntryForm := m.entryForm.(EntryFormModel)
				mEntryForm.Focus()
				m.entryForm = mEntryForm
				return m, func() tea.Msg { return EntryFormMsg{} }
			}
		case "d":
			// delete Todo
			if m.state == tableView {
				mTable := m.table.(TodoTableModel)
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
		case "t":
			if m.state == tableView {
				mTable := m.table.(TodoTableModel)
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
		case "enter":
			if m.state == entryFormView {
				mTable := m.table.(TodoTableModel)
				mEntryForm := m.entryForm.(EntryFormModel)
				title := mEntryForm.Value()
				if title == "" {
					return m, nil
				}
				if err := mTable.AddTodo(title); err != nil {
					// TODO: do something here
				}
				mEntryForm = NewEntryForm()
				mEntryForm.Blur()
				m.entryForm = mEntryForm
				mTable.Focus()
				m.table = mTable
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
		if m.extra != "" {
			return m.table.View() + "\n" + m.extra + "\n"
		}
		return m.table.View()
	case entryFormView:
		return m.entryForm.View()
	default:
		return "Something has gone wrong ..."
	}
}
