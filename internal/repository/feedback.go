package repository

import (
	"context"
	"database/sql"
	"time"
)

type Feedback struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	AppName   string    `json:"app_name"`
	ViewName  string    `json:"view_name"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UserEmail string    `json:"user_email,omitempty"`
}

type FeedbackRepository struct {
	db *sql.DB
}

func NewFeedbackRepository(db *sql.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) Create(ctx context.Context, userID int, appName, viewName, content string) error {
	query := `INSERT INTO feedback (user_id, app_name, view_name, content, status, created_at) VALUES ($1, $2, $3, $4, 'waiting_review', NOW())`
	_, err := r.db.ExecContext(ctx, query, userID, appName, viewName, content)
	return err
}

func (r *FeedbackRepository) List(ctx context.Context, statusFilter string) ([]Feedback, error) {
	query := `
		SELECT f.id, f.user_id, f.app_name, f.view_name, f.content, f.status, f.created_at, u.email
		FROM feedback f
		JOIN users u ON f.user_id = u.id
	`
	var rows *sql.Rows
	var err error

	if statusFilter != "" {
		query += ` WHERE f.status = $1 ORDER BY f.created_at DESC`
		rows, err = r.db.QueryContext(ctx, query, statusFilter)
	} else {
		query += ` ORDER BY f.created_at DESC`
		rows, err = r.db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var f Feedback
		if err := rows.Scan(&f.ID, &f.UserID, &f.AppName, &f.ViewName, &f.Content, &f.Status, &f.CreatedAt, &f.UserEmail); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

func (r *FeedbackRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE feedback SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
