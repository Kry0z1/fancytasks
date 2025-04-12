package database

import (
	"context"
	"database/sql"
	"errors"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/lib/pq"
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
		if err := rows.Scan(&nt.ID, &nt.Title, &nt.Description, &nt.Done, &nt.Owner, &nt.StartsAt, &nt.EndsAt, &nt.Period, &nt.Loop, pq.Array(&nt.Except)); err != nil {
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

func CreateBaseTask(ctx context.Context, task *tasks.BaseTask) error {
	res, err := db.ExecContext(
		ctx,
		"INSERT INTO base_tasks(title, description, done, owner) VALUES ($1, $2, $3, $4)",
		task.Title, task.Description, task.Done, task.Owner,
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	task.ID = int(id)

	return nil
}

func CreateEvent(ctx context.Context, task *tasks.Event) error {
	res, err := db.ExecContext(
		ctx,
		`INSERT INTO events(title, description, done, owner, starts_at, ends_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		task.Title, task.Description, task.Done, task.Owner, task.StartsAt, task.EndsAt,
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	task.ID = int(id)

	return nil
}

func CreateTaskWithDeadline(ctx context.Context, task *tasks.TaskWithDeadline) error {
	res, err := db.ExecContext(
		ctx,
		`INSERT INTO tasks_with_deadline(title, description, done, owner, deadline) 
		VALUES ($1, $2, $3, $4, $5)`,
		task.Title, task.Description, task.Done, task.Owner, task.Deadline,
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	task.ID = int(id)

	return nil
}

func CreateRepeatingTask(ctx context.Context, task *tasks.RepeatingTask) error {
	res, err := db.ExecContext(
		ctx,
		`INSERT INTO repeating_tasks(title, description, done, owner, starts_at, ends_at, period, loop, excepts) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		task.Title, task.Description, task.Done, task.Owner, task.StartsAt, task.EndsAt, task.Period, task.Loop, pq.Array(task.Except),
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	task.ID = int(id)

	return nil
}

// Inserts values from db to task.
// When returned func is called, task values update in db
//
// If error is nil, then returned function should be called, otherwise connection to db will hang
func UpdateBaseTask(ctx context.Context, task *tasks.BaseTask) (func(context.Context) error, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(
		ctx,
		`SELECT title, description, done, owner FROM base_tasks WHERE id = $1 FOR UPDATE`,
		task.ID,
	).Scan(&task.Title, &task.Description, &task.Done, &task.Owner)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return func(ctx context.Context) error {
		defer tx.Rollback()

		_, err := tx.ExecContext(
			ctx,
			`UPDATE base_tasks SET title=$1,description=$2,done=$3,owner=$4
			 WHERE id=$5`,
			task.Title, task.Description, task.Done, task.Owner, task.ID,
		)

		if err != nil {
			return err
		}

		tx.Commit()

		return nil
	}, nil
}

// Inserts values from db to task.
// When returned func is called, task values update in db
//
// If error is nil, then returned function should be called, otherwise connection to db will hang
func UpdateEvent(ctx context.Context, task *tasks.Event) (func(context.Context) error, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(
		ctx,
		`SELECT title, description, done, owner, starts_at, ends_at FROM events WHERE id = $1 FOR UPDATE`,
		task.ID,
	).Scan(&task.Title, &task.Description, &task.Done, &task.Owner, &task.StartsAt, &task.EndsAt)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return func(ctx context.Context) error {
		defer tx.Rollback()

		_, err := tx.ExecContext(
			ctx,
			`UPDATE events SET title=$1,description=$2,done=$3,owner=$4,starts_at=$5,ends_at=$6
			 WHERE id=$7`,
			task.Title, task.Description, task.Done, task.Owner, task.StartsAt, task.EndsAt, task.ID,
		)

		if err != nil {
			return err
		}

		tx.Commit()

		return nil
	}, nil
}

// Inserts values from db to task.
// When returned func is called, task values update in db
//
// If error is nil, then returned function should be called, otherwise connection to db will hang
func UpdateTaskWithDeadline(ctx context.Context, task *tasks.TaskWithDeadline) (func(context.Context) error, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(
		ctx,
		`SELECT title, description, done, owner, deadline FROM tasks_with_deadline WHERE id = $1 FOR UPDATE`,
		task.ID,
	).Scan(&task.Title, &task.Description, &task.Done, &task.Owner, &task.Deadline)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return func(ctx context.Context) error {
		defer tx.Rollback()

		_, err := tx.ExecContext(
			ctx,
			`UPDATE tasks_with_deadline SET title=$1,description=$2,done=$3,owner=$4,deadline=$5
			 WHERE id=$6`,
			task.Title, task.Description, task.Done, task.Owner, task.Deadline, task.ID,
		)

		if err != nil {
			return err
		}

		tx.Commit()

		return nil
	}, nil
}

// Inserts values from db to task.
// When returned func is called, task values update in db
//
// If error is nil, then returned function should be called, otherwise connection to db will hang
func UpdateRepeatingTask(ctx context.Context, task *tasks.RepeatingTask) (func(context.Context) error, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(
		ctx,
		`SELECT title, description, done, owner, starts_at, ends_at, period, loop, excepts 
		FROM repeating_tasks WHERE id = $1 FOR UPDATE`,
		task.ID,
	).Scan(&task.Title, &task.Description, &task.Done, &task.Owner,
		&task.StartsAt, &task.EndsAt, &task.Period, &task.Loop, pq.Array(&task.Except))

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return func(ctx context.Context) error {
		defer tx.Rollback()

		_, err := tx.ExecContext(
			ctx,
			`UPDATE repeating_tasks SET title=$1,description=$2,done=$3,owner=$4,starts_at=$5,
			ends_at=$6,period=$7,loop=$8,excepts=$9
			 WHERE id=$10`,
			task.Title, task.Description, task.Done, task.Owner, task.StartsAt,
			task.EndsAt, task.Period, task.Loop, pq.Array(task.Except), task.ID,
		)

		if err != nil {
			return err
		}

		tx.Commit()

		return nil
	}, nil
}
