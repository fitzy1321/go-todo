package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fitzy1321/go-todo/internal/db"
)

var StartingModel tea.Model

func Start(db *db.AppDB) tea.Model {
	if StartingModel == nil {
		StartingModel = NewTodoTable(db)
	}
	return StartingModel
	// return NewTodoModel(db)
}
