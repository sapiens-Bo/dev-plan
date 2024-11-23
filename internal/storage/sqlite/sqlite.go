package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sapiens-Bo/dev-plan/internal/desk"
	"github.com/sapiens-Bo/dev-plan/internal/task"
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

func (s *Storage) CreateTask(deskID int64, description string, parentTask *int64) (*task.Task, error) {
	const op = "storage.sqlite.CreateTask"
	var parentVal interface{} = nil
	if parentTask != nil {
		parentVal = parentTask
	}

	stmt, err := s.db.Prepare(`
    INSERT INTO tasks(deskID, parent_task_id, description, complited)
    VALUES(?, ?, ?, false)
  `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(deskID, parentVal, description)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return task.New(lastID, deskID, parentTask, description), nil
}

func (s *Storage) GetTasksByDeskID(deskID int64) ([]*task.Task, error) {
	const op = "storage.sqlite.GetTasksByDeskID"
	rows, err := s.db.Query(`
    SELECT id, desk_id, parent_task_id, description, complited
    FROM tasks WHERE desk_id = ?
    ORDER BY id ASC `, deskID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		task := &task.Task{}
		var parentTask sql.NullInt64
		if err := rows.Scan(
			&task.ID,
			&task.DeskID,
			&parentTask,
			&task.Description,
			&task.Complited,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if parentTask.Valid {
			task.ParentTaskID = &parentTask.Int64
		} else {
			task.ParentTaskID = nil
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}
