package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type (
	UserModel struct {
		ID        uuid.UUID
		Name      string
		Email     string
		Password  string
		CreatedAt time.Time
		CreatedBy uuid.UUID
		UpdatedAt pq.NullTime
		UpdatedBy uuid.NullUUID
	}
	UserModelResponse struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Password  string    `json:"-"`
		CreatedAt time.Time `json:"created_at"`
		CreatedBy uuid.UUID `json:"created_by"`
		UpdatedAt time.Time `json:"updated_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
	}
)

// Convert user model into json-friendly formatted response struct (without null data type).
func (user *UserModel) Response() UserModelResponse {
	return UserModelResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt.Time,
		UpdatedBy: user.UpdatedBy.UUID,
	}
}

func (am UserModelResponse) Identifier() uuid.UUID {
	return am.ID
}

func GetAllUser(ctx context.Context, db *sql.DB) ([]UserModel, error) {
	query := `SELECT 
		id,
		name,
		email,
		password,
		created_at,
		created_by,
		updated_at,
		updated_by
	FROM "user"`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userList []UserModel

	for rows.Next() {
		var user UserModel
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.CreatedBy,
			&user.UpdatedAt,
			&user.UpdatedBy,
		)
		if err != nil {
			return userList, errors.Wrap(err, "model/user/list/scan")
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func GetOneUser(ctx context.Context, db *sql.DB, ID uuid.UUID) (UserModel, error) {
	query := `SELECT 
		id,
		name,
		email,
		password,
		created_at,
		created_by,
		updated_at,
		updated_by
	FROM "user" 
	WHERE id=$1`
	var user UserModel
	err := db.QueryRowContext(ctx, query, ID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}
func GetOneUserByEmail(ctx context.Context, db *sql.DB, email string) (UserModel, error) {
	query := `SELECT 
		id,
		name,
		email,
		password,
		created_at,
		created_by,
		updated_at,
		updated_by
	FROM "user" 
	WHERE email=$1`
	var user UserModel
	err := db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (usr UserModel) Insert(ctx context.Context, db *sql.DB) (UserModel, error) {
	query := `INSERT INTO "user" (
		name,
		email,
		password,
		created_by,
		created_at
	) VALUES($1, $2, $3, $4, now()) RETURNING id`
	err := db.QueryRowContext(ctx, query,
		usr.Name,
		usr.Email,
		usr.Password,
		usr.CreatedBy).Scan(
		&usr.ID,
	)
	if err != nil {
		return usr, err
	}
	return usr, nil
}

func (usr UserModel) Update(ctx context.Context, db *sql.DB) error {
	query := `UPDATE "user" SET(
		name,
		email,
		password,
		updated_by,
		updated_at
	)=(
		$1, 
		$2, 
		$3,
		$4,
		now()
	) 
	WHERE id=$5`
	_, err := db.ExecContext(ctx, query,
		usr.Name,
		usr.Email,
		usr.Password,
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
