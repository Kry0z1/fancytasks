package database

import (
	"context"
	"database/sql"
	"errors"

	tasks "github.com/Kry0z1/fancytasks/pkg"
)

func DecorateGetWithTx[T any, V any, E ~[]T | *T](
	ctx context.Context,
	f func(context.Context, *sql.Tx, V) (E, error),
	arg V,
) (E, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := f(ctx, tx, arg)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

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
		`SELECT username, hashed_password FROM users WHERE username = $1`,
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

func GetUserBaseTasks(ctx context.Context, username string) ([]tasks.BaseTask, error) {
	return DecorateGetWithTx[tasks.BaseTask](ctx, GetUserBaseTasksTx, username)
}

// User is responsible for creating and commiting/rollbacking transaction
func GetUserBaseTasksTx(ctx context.Context, tx *sql.Tx, username string) ([]tasks.BaseTask, error) {
	var result []tasks.BaseTask

	rows, err := tx.QueryContext(
		ctx,
		"SELECT * FROM base_tasks WHERE owner = $1",
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var nt tasks.BaseTask
		if err := rows.Scan(&nt.ID, &nt.Title, &nt.Description, &nt.Done, &nt.Owner); err != nil {
			return nil, err
		}
		result = append(result, nt)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func GetUserEvents(ctx context.Context, username string) ([]tasks.Event, error) {
	return DecorateGetWithTx[tasks.Event](ctx, GetUserEventsTx, username)
}

// User is responsible for creating and commiting/rollbacking transaction
func GetUserEventsTx(ctx context.Context, tx *sql.Tx, username string) ([]tasks.Event, error) {
	var result []tasks.Event

	rows, err := tx.QueryContext(
		ctx,
		"SELECT * FROM events WHERE owner = $1",
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var nt tasks.Event
		if err := rows.Scan(&nt.ID, &nt.Title, &nt.Description, &nt.Done, &nt.Owner, &nt.StartsAt, &nt.EndsAt); err != nil {
			return nil, err
		}
		result = append(result, nt)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func GetUserTasksWithDeadline(ctx context.Context, username string) ([]tasks.TaskWithDeadline, error) {
	return DecorateGetWithTx[tasks.TaskWithDeadline](ctx, GetUserTasksWithDeadlineTx, username)
}

// User is responsible for creating and commiting/rollbacking transaction
func GetUserTasksWithDeadlineTx(ctx context.Context, tx *sql.Tx, username string) ([]tasks.TaskWithDeadline, error) {
	var result []tasks.TaskWithDeadline

	rows, err := tx.QueryContext(
		ctx,
		"SELECT * FROM tasks_with_deadline WHERE owner = $1",
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var nt tasks.TaskWithDeadline
		if err := rows.Scan(&nt.ID, &nt.Title, &nt.Description, &nt.Done, &nt.Owner, &nt.Deadline); err != nil {
			return nil, err
		}
		result = append(result, nt)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func GetUserRepeatingTasks(ctx context.Context, username string) ([]tasks.RepeatingTask, error) {
	return DecorateGetWithTx[tasks.RepeatingTask](ctx, GetUserRepeatingTasksTx, username)
}

// User is responsible for creating and commiting/rollbacking transaction
func GetUserRepeatingTasksTx(ctx context.Context, tx *sql.Tx, username string) ([]tasks.RepeatingTask, error) {
	var result []tasks.RepeatingTask

	rows, err := tx.QueryContext(
		ctx,
		"SELECT * FROM repeating_tasks WHERE owner = $1",
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var nt tasks.RepeatingTask
		if err := rows.Scan(&nt.ID, &nt.Title, &nt.Description, &nt.Done, &nt.Owner, &nt.StartsAt, &nt.EndsAt, &nt.Period, &nt.Loop, &nt.Except); err != nil {
			return nil, err
		}
		result = append(result, nt)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

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
		"SELECT username FROM users WHERE username = $1",
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
		"INSERT INTO users(username, hashed_password) VALUES ($1, $2)",
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
