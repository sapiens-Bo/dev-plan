package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sapiens-Bo/dev-plan/internal/desk"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
  CREATE TABLE IF NOT EXISTS desks(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL);

  CREATE TABLE IF NOT EXISTS tasks(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    desk_id INTEGER NOT NULL,
    parent_task_id INTEGER,
    description TEXT NOT NULL,
    complited BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(desk_id) REFERENCES desks(id),
    FOREIGN KEY(parent_task_id) REFERENCES tasks(id));
  `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateDesk(name string) (*desk.Desk, error) {
	const op = "storage.sqlite.CreateDesk"
	stmt, err := s.db.Prepare(`INSERT INTO desks(name) VALUES(?)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	desk := desk.New(lastID, name)

	return desk, nil
}
