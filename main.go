package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fitzy1321/go-todo/internal/db"
	"github.com/fitzy1321/go-todo/internal/tui"
)

func main() {
	adb, err := db.New("todos.db")
	if err != nil {
		fmt.Printf("Bummer about your database: %v", err)
		os.Exit(1)
	}
	if adb == nil {
		fmt.Print("Couldn't get the db! crashing out bruh ...")
		os.Exit(1)
	}
	defer adb.Close()

	p := tea.NewProgram(tui.NewApp(adb), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}
}
