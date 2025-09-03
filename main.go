package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Todo struct {
	ID          uuid.UUID
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time // can be nil
}

type Todos []Todo

func New(title string) *Todo {
	return &Todo{uuid.New(), title, false, time.Now(), nil}
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

type TodoRepo struct {
	db *sql.DB
}

func (t *TodoRepo) CreateTodo(title string) (*Todo, error) {
	ntodo := New(title)
	query := "INSERT INTO todos (id, title, created_at) VALUES (?, ?, ?)"
	_, err := t.db.Exec(
		query,
		ntodo.ID,
		ntodo.Title,
		ntodo.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ntodo, fmt.Errorf("todo item '%s' arleady exists", ntodo.Title)
		}
	}
	return ntodo, err
}

func (tdb *TodoRepo) GetTodos() (todos Todos, err error) {
	query := "SELECT id, title, completed, created_at, completed_at FROM todos ORDER BY created_at DESC"
	rows, err := tdb.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		var idStr string
		var completedAt sql.NullTime

		err := rows.Scan(&idStr, &todo.Title, &todo.Completed, &todo.CreatedAt, &completedAt)
		if err != nil {
			return nil, err
		}

		todo.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			todo.CompletedAt = &completedAt.Time
		}

		todos = append(todos, todo)
	}
	return todos, nil
}

func NewDB() *TodoRepo {
	db, err := sql.Open("sqlite", "todos.db")
	if err != nil {
		if db != nil {
			db.Close()
		}
		log.Fatal(err)
	}

	queries := []string{`
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOL DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME NULL
	)`,
		// tilte unique index
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_todos_title_unique ON todos (title)`,
		// archive table
		`CREATE TABLE IF NOT EXISTS todos_archive (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME NOT NULL,
		completed_at DATETIME NULL,
		archived_at DATETIME NOT NULL
	)`,
		// trigger for archive table
		// prevent updates to archive
		`CREATE TRIGGER IF NOT EXISTS prevent_archive_update
	BEFORE UPDATE ON todos_archive
	BEGIN
		SELECT RAISE(ABORT, 'Archive table is readonly - updates not allowed');
	END`,
		// prevent deletes
		// delete the whole file if you want to delete archives
		`CREATE TRIGGER IF NOT EXISTS prevent_archive_delete
	BEFORE DELETE ON todos_archive
	BEGIN
		SELECT RAISE(ABORT, 'Archive table is readonly - deletes not allowed');
	END`,
	}

	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &TodoRepo{db}
}

func (t *TodoRepo) Close() error {
	return t.db.Close()
}

func main() {
	db := NewDB()
	defer db.Close()
	todos, err := db.GetTodos()
	var titles []string
	for _, t := range todos {
		titles = append(titles, t.Title)
	}
	fmt.Print(titles)
	if err != nil {
		log.Fatal(err)
	}
}
