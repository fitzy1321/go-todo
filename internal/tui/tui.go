package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fitzy1321/go-todo/internal/db"
)

type sessionState int

const (
	tableView sessionState = iota
	addTodoView
)

type AppModel struct {
	db    *db.AppDB
	state sessionState

	entryForm EntryFormModel
	table     TodoTableModel
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
	formMsg  struct{}
	tableMsg struct{}
)

func switchToForm() tea.Msg {
	return formMsg{}
}

func switchToTable() tea.Msg {
	return tableMsg{}
}

func (m AppModel) Init() tea.Cmd { return nil }

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.extra = ""
	switch msg := msg.(type) {
	case formMsg:
		m.state = addTodoView
	case tableMsg:
		m.state = tableView

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "n":
			// Make a new Todo
			if m.state == tableView {
				m.table.Blur()
				m.entryForm.Focus()
				return m, switchToForm
			}
		case "d":
			// delete Todo
			if m.state == tableView {
				todo, err := m.table.SelectedTodo()
				if err != nil {
					m.extra = err.Error()
					return m, nil
				}
				if todo == nil {
					m.extra = "Warn: 'delete cmd' could not find todo object"
					return m, nil
				}
				mt, cmd := m.table.Update(DeleteMsg{id: todo.ID})
				m.table = mt.(TodoTableModel)
				return m, cmd
			}
		case "t":
			if m.state == tableView {
				todo, err := m.table.SelectedTodo()
				if err != nil {
					m.extra = err.Error()
					return m, nil
				}
				if todo == nil {
					m.extra = "Warn: 'toggle cmd' could not find todo object"
					return m, nil
				}
				mt, cmd := m.table.Update(ToggleMsg{id: todo.ID})
				m.table = mt.(TodoTableModel)
				return m, cmd
			}
		case "enter":
			if m.state == addTodoView {
				title := m.entryForm.text.Value()
				if title == "" {
					return m, nil
				}
				if err := m.table.AddTodo(title); err != nil {
					// TODO: do something here
				}
				m.entryForm = NewEntryForm()
				m.entryForm.Blur()
				m.table.Focus()
				return m, switchToTable
			}
		}
	}
	switch m.state {
	case tableView:
		mt, cmd := m.table.Update(msg)
		m.table = mt.(TodoTableModel)
		return m, cmd
	case addTodoView:
		mt, cmd := m.entryForm.Update(msg)
		m.entryForm = mt.(EntryFormModel)
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
	case addTodoView:
		return m.entryForm.View()
	default:
		return "Something has gone wrong ..."
	}
}
