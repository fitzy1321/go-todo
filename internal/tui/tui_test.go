package tui

import (
	"testing"

	"github.com/fitzy1321/go-todo/internal/db"
)

func TestTuiStart(t *testing.T) {
	db, err := db.New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	appModel := NewModel(db)
	if appModel.state != tableView {
		t.Error("Inital appModel.state value is wrong!")
	}

}
