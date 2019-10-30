
	package model

	import (
		"context"
		"database/sql"
		"time"

		"github.com/lib/pq"
		"github.com/satori/go.uuid"
	)

	type (
		UserModel struct {
			ID        uuid.UUID     `json:"id"`
			Name      string        `json:"name"`
			CreatedAt time.Time     `json:"created_at"`
			CreatedBy uuid.UUID     `json:"created_by"`
			UpdatedAt pq.NullTime   `json:"updated_at"`
			UpdatedBy uuid.NullUUID `json:"updated_by"`
		}
	)

	func GetAllUser(ctx context.Context, db *sql.DB) ([]UserModel, error) {
		var userList []UserModel
		query := `SELECT id, name FROM "user"`
		
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return userList, err
		}
		defer rows.Close()

		for rows.Next() {
			var user UserModel
			err := rows.Scan(
				&user.ID,
				&user.Name,
			)
			if err != nil {
				return userList, err
			}
			userList = append(userList, user)
		}
		return userList, nil
	}

	func GetOneUser(ctx context.Context, db *sql.DB, ID uuid.UUID) (UserModel, error) {
		var user UserModel
		query := `SELECT id, name FROM "user" WHERE id=$1`
		err := db.QueryRowContext(ctx, query, ID).Scan(
			&user.ID,
			&user.Name,
		)
		if err != nil {
			return user, err
		}
		return user, nil
	}

	func (usr UserModel) Insert(ctx context.Context, db *sql.DB) (uuid.UUID, error) {
		var id uuid.UUID

		query := `INSERT INTO "user"(name, created_by)VALUES($1, $2) RETURNING id`
		err := db.QueryRowContext(ctx, query,
			usr.Name,
			usr.CreatedBy).Scan(&id)
		if err != nil {
			return id, err
		}
		return id, nil
	}

	func (usr UserModel) Update(ctx context.Context, db *sql.DB) error {
		query := `UPDATE "user" SET(name, updated_by)=($1, $2) WHERE id=$3`
		_, err := db.ExecContext(ctx, query,
			usr.Name,
			usr.UpdatedBy,
			usr.ID)
		if err != nil {
			return err
		}
		return nil
	}

	func DeleteUser(ctx context.Context, db *sql.DB, ID uuid.UUID) error {
		query := `DELETE FROM "user" WHERE id = $1`
		_, err := db.ExecContext(ctx, query, ID)
		if err != nil {
			return err
		}
		return nil
	}
	