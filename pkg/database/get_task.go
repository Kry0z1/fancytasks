package database

import (
	"context"
	"database/sql"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/lib/pq"
)

func GetUserBaseTasks(ctx context.Context, username string) ([]tasks.BaseTask, error) {
	return DecorateGetWithTx[tasks.BaseTask](ctx, GetUserBaseTasksTx, username)
}

// User is responsible for creating and commiting/rollbacking transaction
func GetUserBaseTasksTx(ctx context.Context, tx *sql.Tx, username string) ([]tasks.BaseTask, error) {
	var result []tasks.BaseTask

	rows, err := tx.QueryContext(
		ctx,
		`SELECT 
			id, title, description, done, owner 
		FROM 
			base_tasks 
		WHERE 
			owner = $1`,
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
		`SELECT 
			id, title, description, done, owner, starts_at, ends_at
		FROM 
			events 
		WHERE 
			owner = $1`,
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
		`SELECT 
			id, title, description, done, owner, deadline
		FROM 
			tasks_with_deadline 
		WHERE 
			owner = $1`,
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
		`SELECT 
			id, title, description, done, owner,
			starts_at, ends_at, period, loop, excepts
		FROM 
			repeating_tasks 
		WHERE 
			owner = $1`,
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
