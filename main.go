package main

import (
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fitzy1321/go-todo/internal/db"
	"github.com/fitzy1321/go-todo/internal/todo"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
	db    *db.AppDB
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
		case "down", "j":
		case "enter", " ":
		}
	}
	return m, nil
}

func (m model) View() string {
	return baseStyle.Render(m.table.View())
}

func initModel(db *db.AppDB) model {
	todos, err := db.ListAllTodos()
	if err != nil {
		log.Fatal(err)
	}
	if todos == nil {
		todos = []todo.Todo{}
	}
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Todo", Width: 30},
		{Title: "Completed", Width: 10},
		{Title: "Created At", Width: 20},
		{Title: "Completed At", Width: 20},
	}
	var rows []table.Row
	for _, r_t := range todos {
		var completedAtStr string
		if r_t.CompletedAt != nil {
			completedAtStr = r_t.CompletedAt.String()
		} else {
			completedAtStr = " ~ "
		}

		rows = append(rows, table.Row{
			r_t.ID.String(),
			r_t.Title,
			strconv.FormatBool(r_t.Completed),
			r_t.CreatedAt.String(),
			completedAtStr,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(100),
	)
	return model{t, db}
}

func main() {
	db := db.New()
	defer db.Close()

	if _, err := tea.NewProgram(initModel(db), tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("An error occured: %v", err)
	}
}
