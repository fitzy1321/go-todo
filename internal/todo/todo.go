package todo

import (
	"time"

	"github.com/google/uuid"
)

// Ideas for other statuses
// type status int

// const (
// 	notStarted status = iota
// 	wip
// 	done
// )

type Todo struct {
	ID          uuid.UUID
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time // pointer makes it nullable
}

type Todos []Todo

type TodoArchive struct {
	ArchivedAt time.Time
	Todo
}

func New(title string) Todo {
	return Todo{uuid.New(), title, false, time.Now(), nil}
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

// func (t *Todo) Toggle() {
// 	t.Completed = !t.Completed
// 	if t.Completed {
// 		now := time.Now()
// 		t.CompletedAt = &now
// 	} else {
// 		t.CompletedAt = nil
// 	}
// }
