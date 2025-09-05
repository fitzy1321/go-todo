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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case formMsg:
		m.state = addTodoView
	case tableMsg:
		m.state = tableView

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "d":
			// TODO: delete Todo
		case "n":
			// Make a new Todo
			if m.state == tableView {
				m.table.Blur()
				m.entryForm.Focus()
				return m, switchToForm
			}
		case "t":
			if m.state == tableView {
				if err := m.table.Toggle(); err != nil {

				}
				mT, cmd := m.table.Update(msg)
				m.table = mT.(TodoTableModel)
				return m, cmd
			}
		case "enter":
			switch m.state {
			case tableView:

			case addTodoView:
				title := m.entryForm.text.Value()
				if title == "" {
					return m, nil
				}
				err := m.table.AddTodo(title)
				if err != nil {
					// TODO: do something here
				}
				m.entryForm.Blur()
				m.table.Focus()
				return m, switchToTable
			}
		}
	}
	switch m.state {
	case tableView:
		var t tea.Model
		t, cmd = m.table.Update(msg)
		m.table = t.(TodoTableModel)
		return m, cmd
	case addTodoView:
		var t tea.Model
		t, cmd = m.entryForm.Update(msg)
		m.entryForm = t.(EntryFormModel)
		return m, cmd
	}
	return m, nil
}

func (m AppModel) View() string {
	switch m.state {
	case tableView:
		return m.table.View()
	case addTodoView:
		return m.entryForm.View()
	}
	return "Under Construction"
}
