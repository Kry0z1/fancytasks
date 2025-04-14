package database

import (
	"context"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/lib/pq"
)

func CreateBaseTask(ctx context.Context, task *tasks.BaseTask) error {
	res, err := db.ExecContext(
		ctx,
		`INSERT INTO 
			base_tasks(title, description, done, owner, topic)
		VALUES 
			($1, $2, $3, $4, $5)`,
		task.Title, task.Description, task.Done, task.Owner, task.Topic,
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
		`INSERT INTO 
			events(title, description, done, owner, starts_at, ends_at, topic) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)`,
		task.Title, task.Description, task.Done, task.Owner, task.StartsAt, task.EndsAt, task.Topic,
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
		`INSERT INTO 
			tasks_with_deadline(title, description, done, owner, deadline, topic) 
		VALUES 
			($1, $2, $3, $4, $5, $6)`,
		task.Title, task.Description, task.Done, task.Owner, task.Deadline, task.Topic,
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
		`INSERT INTO 
			repeating_tasks(title, description, done, owner, starts_at, ends_at, period, loop, excepts, topic) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		task.Title, task.Description, task.Done, task.Owner, task.StartsAt, task.EndsAt, task.Period, task.Loop, pq.Array(task.Except), task.Topic,
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	task.ID = int(id)

	return nil
}
