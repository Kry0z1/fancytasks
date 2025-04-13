package database

import (
	"context"
	"database/sql"

	tasks "github.com/Kry0z1/fancytasks/pkg"
)

func GetUser(ctx context.Context, username string) (*tasks.User, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userPass, err := GetUserWithPassword(ctx, username)
	if err != nil {
		return nil, err
	}

	user, err := GetUserWithTasksTx(ctx, tx, username)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.HashedPassword = userPass.HashedPassword

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserWithPassword(ctx context.Context, username string) (*tasks.User, error) {
	var user tasks.User
	err := db.QueryRowContext(
		ctx,
		`SELECT 
			username, hashed_password 
		FROM 
			users 
		WHERE 
			username = $1`,
		username,
	).Scan(&user.Username, &user.HashedPassword)
	return &user, err
}

func GetUserWithTasks(ctx context.Context, username string) (*tasks.User, error) {
	return DecorateGetWithTx[tasks.User](ctx, GetUserWithTasksTx, username)
}

func GetUserWithTasksTx(ctx context.Context, tx *sql.Tx, username string) (*tasks.User, error) {
	var user tasks.User
	var err error

	user.Username = username

	user.BaseTasks, err = GetUserBaseTasksTx(ctx, tx, username)
	if err != nil {
		return nil, err
	}

	user.Events, err = GetUserEventsTx(ctx, tx, username)
	if err != nil {
		return nil, err
	}

	user.TasksWithDeadline, err = GetUserTasksWithDeadlineTx(ctx, tx, username)
	if err != nil {
		return nil, err
	}

	user.RepeatingTasks, err = GetUserRepeatingTasksTx(ctx, tx, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
