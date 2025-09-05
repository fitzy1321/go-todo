package tui

import (
	"errors"
	"log"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fitzy1321/go-todo/internal/db"
	"github.com/fitzy1321/go-todo/internal/todo"
	"github.com/google/uuid"
)

type TodoTableModel struct {
	db       *db.AppDB
	errStr   string
	showHelp bool
	todos    todo.Todos

	table  table.Model
	keyMap TodoTableKeyMap
	help   help.Model
}

type TodoTableKeyMap struct {
	LineUp   key.Binding
	LineDown key.Binding
	New      key.Binding
	Toggle   key.Binding
	Delete   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k TodoTableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.LineUp, k.LineDown, k.New, k.Toggle, k.Delete, k.Help, k.Quit}
}

func (k TodoTableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LineUp, k.LineDown, k.New, k.Toggle, k.Delete},
		{k.Help, k.Quit},
	}
}

func DefaultTodoTableKeyMap() TodoTableKeyMap {
	tBindings := table.DefaultKeyMap()
	return TodoTableKeyMap{
		LineUp:   tBindings.LineUp,
		LineDown: tBindings.LineDown,
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Help: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
	}
}

/* tea.Model Interface: Init, Update, View */

func (m TodoTableModel) Init() tea.Cmd { return nil }

func (m TodoTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.errStr = ""
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ToggleMsg:
		if err := m.Toggle(msg.Id()); err != nil {
			m.errStr = err.Error()
		}
		m.table.Focus()
	case DeleteMsg:
		if err := m.Delete(msg.Id()); err != nil {
			m.errStr = err.Error()
		}
		m.table.Focus()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Help):
			m.showHelp = !m.showHelp
		}
		m.table, cmd = m.table.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		// added offsets, for funzzies
		m.initTable(msg.Width-2, msg.Height-7)
	}
	return m, nil
}

var tableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("252"))

func (m TodoTableModel) View() string {
	var b strings.Builder
	b.WriteString(tableStyle.Render(m.table.View()) + "\n")
	if m.showHelp {
		b.WriteString(m.help.View(m.keyMap) + "\n")
	}
	if m.errStr != "" {
		b.WriteString(errorStyle.Render(m.errStr))
	}
	return b.String()
}

/* Public Functions */

func NewTodoTable(db *db.AppDB) TodoTableModel {
	return TodoTableModel{
		db:     db,
		todos:  todo.Todos{},
		table:  table.New(),
		keyMap: DefaultTodoTableKeyMap(),
		help:   help.New(),
	}
}

func (m *TodoTableModel) AddTodo(title string) {
	mTodo, err := m.db.CreateTodo(strings.TrimSpace(title))
	if err != nil {
		m.errStr = err.Error()
	}
	m.todos = append(m.todos, mTodo)
	sort.Slice(m.todos, func(i, j int) bool {
		return m.todos[i].CreatedAt.Before(m.todos[j].CreatedAt)
	})

	m.table = table.New(
		table.WithColumns(m.todos.Columns()),
		table.WithRows(m.todos.Rows()),
		table.WithFocused(m.table.Focused()),
		table.WithWidth(m.table.Width()),
		table.WithHeight(m.table.Height()),
	)
}

func (m *TodoTableModel) Blur() {
	m.table.Blur()
}

func (m *TodoTableModel) Focus() {
	m.table.Focus()
}

type (
	ZeroTodosError struct{}
	NilRowError    struct{}
)

func (e ZeroTodosError) Error() string {
	return ""
}
func (e NilRowError) Error() string {
	return ""
}

/*
The error returned may be 3 different types:
 1. std `error` ~ most likely parsing UUID error
 2. `ZeroTodosError` ~ when the todos slice is empty
 3. `NilRowError` ~ when the table cursor is out-of-bonds.
*/
func (m *TodoTableModel) SelectedId() (*uuid.UUID, error) {
	if len(m.todos) == 0 {
		return nil, ZeroTodosError{}
	}

	row := m.table.SelectedRow()
	if row == nil {
		return nil, NilRowError{}
	}
	id, err := uuid.Parse(row[0])
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (m *TodoTableModel) Toggle(id uuid.UUID) error {
	var mTodo *todo.Todo
	var index int
	for i, t := range m.todos {
		if t.Id == id {
			mTodo = &t
			index = i
		}
	}
	if mTodo == nil {
		return errors.New("fucking hell, Toggle is broke")
	}

	mTodo.Toggle()

	if err := m.db.UpdateTodo(*mTodo); err != nil {
		return err
	}
	m.todos[index] = *mTodo
	m.table.SetRows(m.todos.Rows())
	return nil
}

func (m *TodoTableModel) Delete(id uuid.UUID) error {
	var mTodo todo.Todo
	var mTodoId = id.String()
	var index *int
	for i, t := range m.todos {
		if t.Id.String() == mTodoId {
			mTodo = t
			index = &i
			break
		}
	}
	if index == nil {
		return errors.New("could not find todo row to delete")
	}
	if err := m.db.DeleteTodo(mTodo); err != nil {
		return err
	}

	if len(m.todos) > 1 {
		m.todos = append(m.todos[:*index], m.todos[*index+1:]...)
		rows := m.table.Rows()
		for i, r := range rows {
			if r[0] == mTodoId {
				index = &i
				break
			}
		}
		m.table.SetRows(append(rows[:*index], rows[*index+1:]...))
		m.table.SetCursor(*index - 1)
	} else {
		m.todos = todo.Todos{}
		m.table.SetRows([]table.Row{})
		m.table.SetCursor(0)
	}

	return nil
}

/* Helper Functions */

func (m *TodoTableModel) initTable(w, h int) {
	var err error

	m.todos, err = m.db.ListAllTodos()
	if err != nil {
		log.Fatal(err)
	} else if m.todos == nil {
		m.todos = todo.NewTodos()
	}

	m.table = table.New(
		table.WithColumns(m.todos.Columns()),
		table.WithRows(m.todos.Rows()),
		table.WithFocused(true),
		table.WithWidth(w),
		table.WithHeight(h),
	)
}
