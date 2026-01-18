package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type OrganizationWithCount struct {
	ID        int
	Name      string
	UserCount int
}

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// List returns all organizations with the count of associated users
func (r *OrganizationRepository) List(ctx context.Context) ([]OrganizationWithCount, error) {
	query := `
		SELECT o.id, o.name, COUNT(u.id) 
		FROM organizations o
		LEFT JOIN users u ON u.organization_name = o.name
		GROUP BY o.id, o.name
		ORDER BY o.name
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []OrganizationWithCount
	for rows.Next() {
		var o OrganizationWithCount
		if err := rows.Scan(&o.ID, &o.Name, &o.UserCount); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, nil
}

// Rename updates an organization's name and all associated users
func (r *OrganizationRepository) Rename(ctx context.Context, oldName, newName string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Update organizations table
	res, err := tx.ExecContext(ctx, "UPDATE organizations SET name = $1 WHERE name = $2", newName, oldName)
	if err != nil {
		return fmt.Errorf("failed to update organization name: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("organization not found: %s", oldName)
	}

	// 2. Update users table
	_, err = tx.ExecContext(ctx, "UPDATE users SET organization_name = $1 WHERE organization_name = $2", newName, oldName)
	if err != nil {
		return fmt.Errorf("failed to update users: %w", err)
	}

	return tx.Commit()
}

// MoveUsers reassigns users from one organization to another
func (r *OrganizationRepository) MoveUsers(ctx context.Context, fromName, toName string) error {
	// Check if target org exists, if not warn? Or just do it?
	// We'll just do the update on users.
	res, err := r.db.ExecContext(ctx, "UPDATE users SET organization_name = $1 WHERE organization_name = $2", toName, fromName)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no users found in organization: %s", fromName)
	}
	return nil
}

// Delete removes an organization from the list (does not delete users, just detaches org info implied?)
// User requirement: "remove organizations from the list". 
// Implementation: Delete row from organizations. Users keep the string but it's no longer a "listed" org.
func (r *OrganizationRepository) Delete(ctx context.Context, name string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM organizations WHERE name = $1", name)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("organization not found: %s", name)
	}
	return nil
}
