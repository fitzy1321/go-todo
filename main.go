package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fitzy1321/go-todo/internal/db"
	"github.com/fitzy1321/go-todo/internal/tui"
)

// func (m Model) Init() tea.Cmd { return m.spinner.Tick }
var Model *tui.TodoTableModel

func main() {
	adb := db.New()
	defer adb.Close()

	Model = tui.NewTodoModel(adb)

	p := tea.NewProgram(Model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("An error occured: %v", err)
	}
}
