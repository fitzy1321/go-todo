package tui

import (
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

type TodoTableModel struct {
	db *db.AppDB
	// cursorRow int // x value
	// cursorCol int // y value
	// loaded bool
	// spinner spinner.Model
	table table.Model
	todos todo.Todos
	extra string // extra text to print
}

func NewTodoTable(db *db.AppDB) *TodoTableModel {
	return &TodoTableModel{db: db, table: table.New()}
}

/* tea.Model Interface: Init, Update, View */
func (m TodoTableModel) Init() tea.Cmd { return nil }

func (m TodoTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.extra = ""
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		mstStr := msg.String()
		// m.extra = "Key stroke: " + mstStr
		switch mstStr {
		case "ctrl+c", "q":
			cmd = tea.Quit
		case "up", "k":
			// if m.cursorRow > 0 {
			// 	m.cursorRow--
			// }
			m.table, cmd = m.table.Update(msg)
		case "down", "j":
			// if m.cursorRow < len(m.table.Rows()) {
			// 	m.cursorRow--
			// }
			m.table, cmd = m.table.Update(msg)
		case "enter", " ":
			return NewTodoInput(), nil
		case "n":
			// TODO: Create New Todo.
			// TODO: Need a form input
			return m, cmd
		}
		return m, cmd

	case tea.WindowSizeMsg:
		w_off, h_off := 2, 10
		m.initTable(msg.Width-w_off, msg.Height-h_off)
	}
	return m, nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("252"))

func (m TodoTableModel) View() string {
	res := baseStyle.Render(m.table.View())
	// if m.extra != "" {
	// 	res += fmt.Sprintf(" \n\n%s", m.extra)
	// }

	return res
}

/* Public Functions */

func AddRow(row table.Row, t *table.Model) *table.Model {
	newt := table.New(
		table.WithColumns(t.Columns()),
		table.WithRows(append(t.Rows(), row)),
		table.WithFocused(true),
		table.WithWidth(t.Width()),
		table.WithHeight(t.Height()),
	)
	return &newt
}

/* Helper Functions */
func (m *TodoTableModel) initTable(width, height int) {
	var err error
	m.todos, err = m.db.ListAllTodos()
	if err != nil {
		log.Fatal(err)
	}
	if m.todos == nil {
		m.todos = todo.NewTodos()
	}

	var titleLen []int = make([]int, len(m.todos))
	for _, item := range m.todos.GetTitles() {
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
		table.WithWidth(width),
		table.WithHeight(height),
	)
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

func createRows(todos todo.Todos) []table.Row {
	var rows []table.Row
	for _, t := range todos {
		rows = append(rows, createRow(t))
	}
	return rows
}
