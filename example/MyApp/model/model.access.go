package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

type (
	AccessModel struct {
		ID        uuid.UUID     `json:"id"`
		Name      string        `json:"name"`
		CreatedAt time.Time     `json:"created_at"`
		CreatedBy uuid.UUID     `json:"created_by"`
		UpdatedAt pq.NullTime   `json:"updated_at"`
		UpdatedBy uuid.NullUUID `json:"updated_by"`
	}
)

func GetAllAccess(ctx context.Context, db *sql.DB) ([]AccessModel, error) {
	var accessList []AccessModel
	query := "SELECT id, name FROM access"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return accessList, err
	}
	defer rows.Close()
	for rows.Next() {
		var access AccessModel
		err := rows.Scan(
			&access.ID,
			&access.Name,
		)
		if err != nil {
			return accessList, err
		}
		accessList = append(accessList, access)
	}
	return accessList, nil
}

func GetOneAccess(ctx context.Context, db *sql.DB, ID uuid.UUID) (AccessModel, error) {
	var access AccessModel
	query := "SELECT id, name FROM access WHERE id= $1"
	err := db.QueryRowContext(ctx, query, ID).Scan(
		&access.ID,
		&access.Name,
	)
	if err != nil {
		return access, err
	}
	return access, nil
}

func (access AccessModel) Insert(ctx context.Context, db *sql.DB) (uuid.UUID, error) {
	var id uuid.UUID
	query := "INSERT INTO access (name, created_by) VALUES ($1, $2) RETURNING id"
	err := db.QueryRowContext(ctx, query, access.Name, access.CreatedBy).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (access AccessModel) Update(ctx context.Context, db *sql.DB) error {
	query := "UPDATE access SET (name, updated_by) = ($1, $2) WHERE id = $3"
	_, err := db.ExecContext(ctx, query, access.Name, access.UpdatedBy, access.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAccess(ctx context.Context, db *sql.DB, ID uuid.UUID) error {
	query := "DELETE FROM access WHERE id = $1"
	_, err := db.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}
	return nil
}
