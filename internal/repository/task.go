package repository

import (
	"context"
	"database/sql"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	IsSelected  bool      `json:"is_selected"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) List(ctx context.Context) ([]Task, error) {
	query := `SELECT id, description, is_completed, is_selected, created_at, updated_at FROM tasks ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Description, &t.IsCompleted, &t.IsSelected, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// SyncUpsert inserts or updates a task based on description.
func (r *TaskRepository) SyncUpsert(ctx context.Context, description string, isCompleted bool) error {
	query := `
		INSERT INTO tasks (description, is_completed, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (description) 
		DO UPDATE SET is_completed = $2, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, description, isCompleted)
	return err
}

func (r *TaskRepository) ToggleSelection(ctx context.Context, id int) error {
	query := `UPDATE tasks SET is_selected = NOT is_selected, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
