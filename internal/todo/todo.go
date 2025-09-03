package todo

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time // pointer makes it nullable
}

type TodoArchive struct {
	ArchivedAt time.Time
	Todo
}

func New(title string) *Todo {
	return &Todo{uuid.New(), title, false, time.Now(), nil}
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
