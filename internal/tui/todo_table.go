package tui

import (
	"errors"
	"log"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fitzy1321/go-todo/internal/db"
	"github.com/fitzy1321/go-todo/internal/todo"
	"github.com/google/uuid"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("252"))

type TodoTableModel struct {
	db     *db.AppDB
	errStr string

	table table.Model
	todos todo.Todos
}

type (
	ToggleMsg struct{ id uuid.UUID }
	DeleteMsg struct{ id uuid.UUID }
)

func (t *ToggleMsg) Id() uuid.UUID {
	return t.id
}

func (t *DeleteMsg) Id() uuid.UUID {
	return t.id
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
		m.table, cmd = m.table.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		w_off, h_off := 2, 10
		m.initTable(msg.Width-w_off, msg.Height-h_off)
	}
	return m, nil
}

func (m TodoTableModel) View() string {
	res := baseStyle.Render(m.table.View())
	if m.errStr != "" {
		res += "\n" + m.errStr
	}
	return res
}

/* Public Functions */

func NewTodoTable(db *db.AppDB) TodoTableModel {
	return TodoTableModel{db: db, table: table.New(), todos: todo.Todos{}}
}

func (m *TodoTableModel) AddTodo(title string) error {
	t, err := m.db.CreateTodo(title)
	if err != nil {
		return err
	}
	m.todos = append(m.todos, t)

	m.table = table.New(
		table.WithColumns(m.table.Columns()),
		table.WithRows(append(m.table.Rows(), t.IntoRow())),
		table.WithFocused(m.table.Focused()),
		table.WithWidth(m.table.Width()),
		table.WithHeight(m.table.Height()),
	)
	return nil
}

func (m *TodoTableModel) Blur() {
	m.table.Blur()
}

func (m *TodoTableModel) Focus() {
	m.table.Focus()
}

func (m *TodoTableModel) Toggle(id uuid.UUID) error {
	var todo *todo.Todo
	var index int
	for i, t := range m.todos {
		if t.ID == id {
			todo = &t
			index = i
		}
	}
	if todo == nil {
		return errors.New("fucking hell, Toggle is broke")
	}

	todo.Toggle()

	if err := m.db.UpdateTodo(*todo); err != nil {
		return err
	}
	m.todos[index] = *todo
	m.table.SetRows(createRows(m.todos))
	return nil
}

func (m *TodoTableModel) SelectedTodo() (*todo.Todo, error) {
	if len(m.todos) == 0 {
		return nil, nil
	}
	row := m.table.SelectedRow()
	id := uuid.MustParse(row[0])
	for _, item := range m.todos {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, errors.New("could not find todo")
}

func (m *TodoTableModel) Delete(id uuid.UUID) error {
	var mTodo todo.Todo
	var index *int
	for i, t := range m.todos {
		if t.ID == id {
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
		m.todos = append(m.todos[:*index], m.todos[(*index)+1:]...)
		rows := m.table.Rows()
		for i, r := range rows {
			if r[0] == mTodo.ID.String() {
				index = &i
				break
			}
		}
		m.table.SetRows(append(rows[:*index], rows[*index+1:]...))
	} else {
		m.todos = todo.Todos{}
		m.table.SetRows([]table.Row{})
	}

	return nil
}

/* Helper Functions */

func (m *TodoTableModel) initTable(w, h int) {
	var err error
	var titleW int = 4

	m.todos, err = m.db.ListAllTodos()
	if err != nil {
		log.Fatal(err)
	}
	if m.todos == nil {
		m.todos = todo.NewTodos()
	} else if len(m.todos) > 0 {
		var titleLens []int
		for _, item := range m.todos.GetTitles() {
			titleLens = append(titleLens, len(item))
		}
		titleW = slices.Max(titleLens)
	}

	uuidW := len(uuid.New().String())
	timeW := len(time.Now().String())

	columns := []table.Column{
		{Title: "ID", Width: uuidW},
		{Title: "Todo", Width: titleW},
		{Title: "Completed", Width: 9},
		{Title: "Created At", Width: timeW},
		{Title: "Completed At", Width: timeW},
	}

	m.table = table.New(
		table.WithColumns(columns),
		table.WithRows(createRows(m.todos)),
		table.WithFocused(true),
		table.WithWidth(w),
		table.WithHeight(h),
	)
}

func createRows(ts todo.Todos) []table.Row {
	var rows []table.Row
	for _, t := range ts {
		rows = append(rows, t.IntoRow())
	}
	return rows
}
