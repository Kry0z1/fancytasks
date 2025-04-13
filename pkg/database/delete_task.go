package database

import (
	"context"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/lib/pq"
)

func DeleteBaseTask(ctx context.Context, task *tasks.BaseTask) error {
	return db.QueryRowContext(
		ctx,
		`DELETE FROM
			base_tasks
		WHERE
			id = $1
		RETURNING
			id, title, description, done, owner`,
		task.ID,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.Owner)
}

func DeleteEvent(ctx context.Context, task *tasks.Event) error {
	return db.QueryRowContext(
		ctx,
		`DELETE FROM
			events
		WHERE
			id = $1
		RETURNING
			id, title, description, done, owner, starts_at, ends_at`,
		task.ID,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.Owner,
		&task.StartsAt, &task.EndsAt)
}

func DeleteTaskWithDeadline(ctx context.Context, task *tasks.TaskWithDeadline) error {
	return db.QueryRowContext(
		ctx,
		`DELETE FROM
			tasks_with_deadline
		WHERE
			id = $1
		RETURNING
			id, title, description, done, owner, deadline`,
		task.ID,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.Owner, &task.Deadline)
}

func DeleteRepeatingTask(ctx context.Context, task *tasks.RepeatingTask) error {
	return db.QueryRowContext(
		ctx,
		`DELETE FROM
			repeating_tasks
		WHERE
			id = $1
		RETURNING
			id, title, description, done, owner, starts_at, ends_at, period, loop, excepts`,
		task.ID,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.Owner,
		&task.StartsAt, &task.EndsAt, &task.Period, &task.Loop, pq.Array(&task.Except))
}
