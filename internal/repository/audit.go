package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID         int             `json:"id"`
	UserID     int             `json:"user_id"`
	Action     string          `json:"action"`
	EntityType string          `json:"entity_type"`
	EntityID   int             `json:"entity_id"`
	Details    json.RawMessage `json:"details"`
	IPAddress  string          `json:"ip_address"`
	CreatedAt  time.Time       `json:"created_at"`
	UserEmail  string          `json:"user_email,omitempty"`
}

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Log(ctx context.Context, userID int, action, entityType string, entityID int, details interface{}, ip string) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return err
	}
	query := `INSERT INTO audit_logs (user_id, action, entity_type, entity_id, details, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW())`
	_, err = r.db.ExecContext(ctx, query, userID, action, entityType, entityID, detailsJSON, ip)
	return err
}

func (r *AuditRepository) List(ctx context.Context, limit int) ([]AuditLog, error) {
	query := `
		SELECT a.id, a.user_id, a.action, a.entity_type, a.entity_id, a.details, a.ip_address, a.created_at, u.email
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var l AuditLog
		var details []byte
		var ip sql.NullString
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.EntityType, &l.EntityID, &details, &ip, &l.CreatedAt, &l.UserEmail); err != nil {
			return nil, err
		}
		l.Details = json.RawMessage(details)
		l.IPAddress = ip.String
		logs = append(logs, l)
	}
	return logs, nil
}
