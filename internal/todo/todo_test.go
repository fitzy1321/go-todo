package todo_test

import (
	"slices"
	"testing"

	"github.com/fitzy1321/go-todo/internal/todo"
	"github.com/google/uuid"
)

func TestGetTitle(t *testing.T) {
	titles := []string{"Foo", "Bar", "Baz"}
	var todos todo.Todos
	for _, item := range titles {
		todos = append(todos, todo.Todo{
			Id:    uuid.New(),
			Title: item,
		})
	}

	actual := todos.GetTitles()
	t.Log(actual)
	if !slices.Equal(titles, actual) {
		t.FailNow()
	}
}
