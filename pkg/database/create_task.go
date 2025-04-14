package database

import (
	"context"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/lib/pq"
)

func CreateBaseTask(ctx context.Context, task *tasks.BaseTask) error {
	return db.QueryRowContext(
		ctx,
		`INSERT INTO 
			base_tasks(title, description, done, owner, topic)
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING
			id`,
		task.Title, task.Description, task.Done, task.Owner, task.Topic,
	).Scan(&task.ID)
}

func CreateEvent(ctx context.Context, task *tasks.Event) error {
	return db.QueryRowContext(
		ctx,
		`INSERT INTO 
			events(title, description, done, owner, starts_at, ends_at, topic) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id`,
		task.Title, task.Description, task.Done, task.Owner, task.StartsAt, task.EndsAt, task.Topic,
	).Scan(&task.ID)
}

func CreateTaskWithDeadline(ctx context.Context, task *tasks.TaskWithDeadline) error {
	return db.QueryRowContext(
		ctx,
		`INSERT INTO 
			tasks_with_deadline(title, description, done, owner, deadline, topic) 
		VALUES 
			($1, $2, $3, $4, $5, $6)
		RETURNING
			id`,
		task.Title, task.Description, task.Done, task.Owner, task.Deadline, task.Topic,
	).Scan(&task.ID)
}

func CreateRepeatingTask(ctx context.Context, task *tasks.RepeatingTask) error {
	return db.QueryRowContext(
		ctx,
		`INSERT INTO 
			repeating_tasks(title, description, done, owner, starts_at, ends_at, period, loop, excepts, topic) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
			id`,
		task.Title, task.Description, task.Done, task.Owner, task.StartsAt, task.EndsAt, task.Period, task.Loop, pq.Array(task.Except), task.Topic,
	).Scan(&task.ID)
}
