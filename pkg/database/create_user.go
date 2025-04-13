package database

import (
	"context"
	"database/sql"
	"errors"

	tasks "github.com/Kry0z1/fancytasks/pkg"
)

var ErrUserExists = errors.New("User with such username already exists")

func CreateUser(ctx context.Context, username string, password string, h tasks.Hasher) (*tasks.User, error) {
	var user tasks.User

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(
		ctx,
		`SELECT 
			username 
		FROM 
			users 
		WHERE 
			username = $1`,
		username,
	).Scan(&user.Username)
	if err == nil {
		return nil, ErrUserExists
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	hashedPassword, err := h.HashPassword(password)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO 
			users(username, hashed_password) 
		VALUES 
			($1, $2)`,
		username, hashedPassword,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	user.Username = username
	user.HashedPassword = hashedPassword

	return &user, nil
}
