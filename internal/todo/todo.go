package todo

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type Todos []Todo

type TodoArchive struct {
	ArchivedAt time.Time
	Todo
}

func New(title string) Todo {
	return Todo{uuid.New(), title, false, time.Now(), nil}
}

func (t *Todo) Toggle() {
	t.Completed = !t.Completed
	if t.Completed {
		now := time.Now()
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
}

func (t Todo) IntoRow() table.Row {
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

func NewTodos() Todos {
	return Todos{}
}

func (t *Todos) GetTitles() []string {
	var titles []string
	for _, item := range *t {
		titles = append(titles, item.Title)
	}
	return titles
}
