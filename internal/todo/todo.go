package todo

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/google/uuid"
)

type Todo struct {
	Id          uuid.UUID
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
	return Todo{
		Id:          uuid.New(),
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
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

func (t Todo) Row() table.Row {
	completedAtStr := " ~ "
	if t.CompletedAt != nil {
		completedAtStr = t.CompletedAt.String()
	}
	return table.Row{
		t.Id.String(),
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

func (t *Todos) Rows() []table.Row {
	var rows []table.Row
	for _, todo := range *t {
		rows = append(rows, todo.Row())
	}
	return rows
}

func (t *Todos) Columns() []table.Column {
	var bytes [36]byte
	uuidW := len(string(bytes[:]))
	timeW := len(strings.TrimSpace(time.Now().String()))

	var titleW int = 4
	if len(*t) > 0 {
		var titleLens []int
		for _, item := range t.GetTitles() {
			titleLens = append(titleLens, len(item))
		}
		titleW = slices.Max(titleLens)
	}

	return []table.Column{
		{Title: "ID", Width: uuidW},
		{Title: "Todo", Width: titleW},
		{Title: "Completed", Width: 9},
		{Title: "Created At", Width: timeW},
		{Title: "Completed At", Width: timeW},
	}
}
