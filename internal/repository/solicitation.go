package repository

import (
	"bd_bot/internal/scraper"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type SolicitationRepository struct {
	db *sql.DB
}

func NewSolicitationRepository(db *sql.DB) *SolicitationRepository {
	return &SolicitationRepository{db: db}
}

// Upsert inserts a solicitation or updates it if source_id already exists
func (r *SolicitationRepository) Upsert(ctx context.Context, sol scraper.Solicitation) error {
	rawData, err := json.Marshal(sol.RawData)
	if err != nil {
		return fmt.Errorf("error marshalling raw data: %w", err)
	}

	docsData, err := json.Marshal(sol.Documents)
	if err != nil {
		return fmt.Errorf("error marshalling documents: %w", err)
	}

	query := `
		INSERT INTO solicitations (source_id, title, description, agency, due_date, url, raw_data, documents, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (source_id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			agency = EXCLUDED.agency,
			due_date = EXCLUDED.due_date,
			url = EXCLUDED.url,
			raw_data = EXCLUDED.raw_data,
			documents = EXCLUDED.documents,
			updated_at = NOW();
	`

	// Handle zero time for due_date
	var dueDate interface{}
	if !sol.DueDate.IsZero() {
		dueDate = sol.DueDate
	}

	_, err = r.db.ExecContext(ctx, query,
		sol.SourceID,
		sol.Title,
		sol.Description,
		sol.Agency,
		dueDate,
		sol.URL,
		rawData,
		docsData,
	)

	return err
}

// List retrieves all solicitations from the database
func (r *SolicitationRepository) List(ctx context.Context) ([]scraper.Solicitation, error) {
	query := `
		SELECT source_id, title, description, agency, due_date, url, raw_data, documents
		FROM solicitations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var solicitations []scraper.Solicitation
	for rows.Next() {
		var sol scraper.Solicitation
		var rawData []byte
		var docsData []byte
		var dueDate sql.NullTime

		if err := rows.Scan(
			&sol.SourceID,
			&sol.Title,
			&sol.Description,
			&sol.Agency,
			&dueDate,
			&sol.URL,
			&rawData,
			&docsData,
		); err != nil {
			return nil, err
		}

		if dueDate.Valid {
			sol.DueDate = dueDate.Time
		}

		if err := json.Unmarshal(rawData, &sol.RawData); err != nil {
			return nil, fmt.Errorf("error unmarshalling raw data: %w", err)
		}

		if len(docsData) > 0 {
			if err := json.Unmarshal(docsData, &sol.Documents); err != nil {
				return nil, fmt.Errorf("error unmarshalling documents: %w", err)
			}
		}

		solicitations = append(solicitations, sol)
	}

	return solicitations, rows.Err()
}
