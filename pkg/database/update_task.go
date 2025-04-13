package database

import (
	"context"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/lib/pq"
)

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
		`SELECT 
			title, description, done, owner 
		FROM 
			base_tasks 
		WHERE 
			id = $1 
		FOR UPDATE`,
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
			`UPDATE 
				base_tasks 
			SET 
				title=$1,description=$2,done=$3,owner=$4
			WHERE 
				id=$5`,
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
		`SELECT 
			title, description, done, owner, starts_at, ends_at 
		FROM 
			events 
		WHERE 
			id = $1 
		FOR UPDATE`,
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
			`UPDATE 
				events 
			SET 
				title=$1,description=$2,done=$3,owner=$4,starts_at=$5,ends_at=$6
			WHERE 
				id=$7`,
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
		`SELECT 
			title, description, done, owner, deadline 
		FROM 
			tasks_with_deadline 
		WHERE 
			id = $1 
		FOR UPDATE`,
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
			`UPDATE 
				tasks_with_deadline 
			SET 
				title=$1,description=$2,done=$3,owner=$4,deadline=$5
			WHERE 
				id=$6`,
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
		`SELECT 
			title, description, done, owner, starts_at, ends_at, period, loop, excepts 
		FROM 
			repeating_tasks 
		WHERE 
			id = $1 
		FOR UPDATE`,
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
			`UPDATE 
				repeating_tasks
			SET 
				title=$1,description=$2,done=$3,owner=$4,starts_at=$5,
				ends_at=$6,period=$7,loop=$8,excepts=$9
			WHERE 
				id=$10`,
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
