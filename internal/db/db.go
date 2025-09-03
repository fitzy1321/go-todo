package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/fitzy1321/go-todo/internal/todo"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type AppDB struct {
	db *sql.DB
}

func New() *AppDB {
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
		// Unique index on title
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_todos_title_unique ON todos (title)`,
		// Archive table
		`CREATE TABLE IF NOT EXISTS todos_archive (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME NOT NULL,
		completed_at DATETIME NULL,
		archived_at DATETIME NOT NULL
	)`,
		// Trigger for archive table
		// Prevent update trigger
		`CREATE TRIGGER IF NOT EXISTS prevent_archive_update
	BEFORE UPDATE ON todos_archive
	BEGIN
		SELECT RAISE(ABORT, 'Archive table is readonly - updates not allowed');
	END`,
		// Prevent delete trigger. Delete '*.db' file to delete archives
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

	return &AppDB{db}
}

func (a *AppDB) Close() error {
	return a.db.Close()
}

func (a *AppDB) CreateTodo(title string) (*todo.Todo, error) {
	ntodo := todo.New(title)
	query := "INSERT INTO todos (id, title, created_at) VALUES (?, ?, ?)"
	_, err := a.db.Exec(
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

func (a *AppDB) ListAllTodos() (todos []todo.Todo, err error) {
	query := "SELECT id, title, completed, created_at, completed_at FROM todos ORDER BY created_at DESC"
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo todo.Todo
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

func (a *AppDB) ListAllArchives() (archive []todo.TodoArchive, err error) {
	query := "SELECT id, title, completed, created_at, completed_at, archived_at FROM todos ORDER BY archived_at DESC"
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo todo.TodoArchive
		var idStr string
		var completedAt sql.NullTime

		err := rows.Scan(&idStr, &todo.Title, &todo.Completed, &todo.CreatedAt, &completedAt, &todo.ArchivedAt)
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

		archive = append(archive, todo)
	}
	return archive, nil
}
