package tui

import (
	"errors"
	"log"
	"slices"
	"strconv"
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
	db    *db.AppDB
	table table.Model
	todos todo.Todos
}

/* tea.Model Interface: Init, Update, View */

func (m TodoTableModel) Init() tea.Cmd { return nil }

func (m TodoTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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
	return baseStyle.Render(m.table.View())
}

/* Public Functions */

func NewTodoTable(db *db.AppDB) TodoTableModel {
	return TodoTableModel{db: db, table: table.New()}
}

func (m *TodoTableModel) AddTodo(title string) error {
	t, err := m.db.CreateTodo(title)
	if err != nil {
		return err
	}
	m.todos = append(m.todos, t)
	newRow := createRow(t)
	m.table = table.New(
		table.WithColumns(m.table.Columns()),
		table.WithRows(append(m.table.Rows(), newRow)),
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

func (m *TodoTableModel) Toggle() error {
	row := m.table.SelectedRow()
	id := row[0]
	var todo *todo.Todo
	for _, t := range m.todos {
		if t.ID == uuid.MustParse(id) {
			todo = &t
		}
	}
	if todo == nil {
		return errors.New("fucking hell, Toggle is broke")
	}
	if !todo.Completed {
		nowT := time.Now()
		todo.CompletedAt = &nowT
	} else {
		todo.CompletedAt = nil
	}
	todo.Completed = !todo.Completed

	if err := m.db.UpdateTodo(*todo); err != nil {
		return err
	}
	m.table.SetRows(createRows(m.todos))
	return nil
}

/* Helper Functions */

func (m *TodoTableModel) initTable(w, h int) {
	todos, err := m.db.ListAllTodos()
	if err != nil {
		log.Fatal(err)
	}
	if todos == nil {
		todos = todo.NewTodos()
	}
	m.todos = todos
	var titleLen []int = make([]int, len(todos))
	for _, item := range todos.GetTitles() {
		titleLen = append(titleLen, len(item))
	}
	uuidW := len(uuid.New().String())
	titleW := slices.Max(titleLen)
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
		rows = append(rows, createRow(t))
	}
	return rows
}
func createRow(t todo.Todo) table.Row {
	completedAtStr := " ~ "
	if t.CompletedAt != nil {
		completedAtStr = t.CompletedAt.String()
	}
	return table.Row{
		t.ID.String(),
		t.Title,
		strconv.FormatBool(t.Completed),
		t.CreatedAt.String(),
		completedAtStr,
	}
}
