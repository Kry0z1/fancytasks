package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Kry0z1/fancytasks/internal/middleware"
	"github.com/Kry0z1/fancytasks/internal/middleware/auth"
	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/database"
)

func UpdateTask(w http.ResponseWriter, r *http.Request) error {
	var (
		id       int
		err      error
		update   func() (func(context.Context) error, error)
		callback func(context.Context) error
		parse    func() error
		send     func() error
		baseTask *tasks.BaseTask
	)

	user := auth.ContextUser(r.Context())
	if user == nil {
		return middleware.HTTPError{
			Err:     nil,
			Message: "Unauthorized",
			Code:    http.StatusUnauthorized,
		}
	}

	if err = r.ParseForm(); err != nil {
		return err
	}

	if id, err = strconv.Atoi(r.Form.Get("id")); err != nil {
		return middleware.HTTPError{
			Err:     nil,
			Message: "Invalid id",
			Code:    http.StatusBadRequest,
		}
	}

	taskType := r.Form.Get("tasktype")
	dctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(time.Second))
	defer cancel()

	switch taskType {
	case "":
		return middleware.HTTPError{
			Err:     nil,
			Message: "Task type not found in form",
			Code:    http.StatusNotFound,
		}
	case "basetask":
		var task tasks.BaseTask
		baseTask = &task
		update = func() (func(context.Context) error, error) { return database.UpdateBaseTask(dctx, &task) }
		parse = func() error { return nil }
		send = func() error { return json.NewEncoder(w).Encode(&task) }
	case "event":
		var task tasks.Event
		baseTask = &task.BaseTask
		update = func() (func(context.Context) error, error) { return database.UpdateEvent(dctx, &task) }
		parse = func() error { return parseEvent(r, &task) }
		send = func() error { return json.NewEncoder(w).Encode(&task) }
	case "deadline":
		var task tasks.TaskWithDeadline
		baseTask = &task.BaseTask
		update = func() (func(context.Context) error, error) { return database.UpdateTaskWithDeadline(dctx, &task) }
		parse = func() error { return parseDeadline(r, &task) }
		send = func() error { return json.NewEncoder(w).Encode(&task) }
	case "repeat":
		var task tasks.RepeatingTask
		baseTask = &task.BaseTask
		update = func() (func(context.Context) error, error) { return database.UpdateRepeatingTask(dctx, &task) }
		parse = func() error { return parseRepeatingTask(r, &task) }
		send = func() error { return json.NewEncoder(w).Encode(&task) }
	default:
		return middleware.HTTPError{
			Err:     nil,
			Message: "Invalid task type",
			Code:    http.StatusBadRequest,
		}
	}

	baseTask.ID = id
	callback, err = update()
	if err == sql.ErrNoRows {
		return middleware.HTTPError{
			Err:     nil,
			Message: "Task with such id not found",
			Code:    http.StatusNotFound,
		}
	}
	if err != nil {
		return err
	}

	if baseTask.Owner != user.Username {
		callback(dctx)
		return middleware.HTTPError{
			Err:     nil,
			Message: "Cannot update tasks of other users",
			Code:    http.StatusUnauthorized,
		}
	}

	parseBaseTask(r, baseTask)
	if err = parse(); err != nil {
		callback(dctx)
		return err
	}

	if err = callback(dctx); err != nil {
		return err
	}
	return send()
}

func parseBaseTask(r *http.Request, task *tasks.BaseTask) {
	if r.Form.Has("title") && r.Form.Get("title") != "" {
		task.Title = r.Form.Get("title")
	}

	if r.Form.Has("description") {
		task.Description = r.Form.Get("description")
	}

	if r.Form.Has("done") {
		task.Done = r.Form.Get("done") == "true"
	}

	if r.Form.Has("topic") {
		task.Topic = r.Form.Get("topic")
	}
}

func parseEvent(r *http.Request, t *tasks.Event) error {
	startsUnix := t.StartsAt.Unix()
	endsUnix := t.EndsAt.Unix()
	var err error

	if r.Form.Has("starts_at") {
		if startsUnix, err = strconv.ParseInt(r.Form.Get("starts_at"), 10, 0); err != nil {
			return middleware.HTTPError{
				Err:     err,
				Message: "Invalid starts_at",
				Code:    http.StatusBadRequest,
			}
		}
	}

	if r.Form.Has("ends_at") {
		if endsUnix, err = strconv.ParseInt(r.Form.Get("ends_at"), 10, 0); err != nil {
			return middleware.HTTPError{
				Err:     err,
				Message: "Invalid ends_at",
				Code:    http.StatusBadRequest,
			}
		}
	}

	if endsUnix < startsUnix {
		return middleware.HTTPError{
			Err:     err,
			Message: "End as earlier than start",
			Code:    http.StatusBadRequest,
		}
	}

	if endsUnix < 0 || startsUnix < 0 {
		return middleware.HTTPError{
			Err:     err,
			Message: "Timestamps cannot be less than 0",
			Code:    http.StatusBadRequest,
		}
	}

	t.StartsAt = time.Unix(startsUnix, 0)
	t.EndsAt = time.Unix(endsUnix, 0)

	return nil
}

func parseDeadline(r *http.Request, t *tasks.TaskWithDeadline) error {
	deadline := t.Deadline.Unix()
	var err error

	if r.Form.Has("deadline") {
		if deadline, err = strconv.ParseInt(r.Form.Get("deadline"), 10, 0); err != nil {
			return middleware.HTTPError{
				Err:     err,
				Message: "Invalid deadline",
				Code:    http.StatusBadRequest,
			}
		}
	}

	if deadline < 0 {
		return middleware.HTTPError{
			Err:     err,
			Message: "Timestamps cannot be less than 0",
			Code:    http.StatusBadRequest,
		}
	}

	t.Deadline = time.Unix(deadline, 0)

	return nil
}

func parseRepeatingTask(r *http.Request, t *tasks.RepeatingTask) error {
	var (
		startsUnix = t.StartsAt.Unix()
		endsUnix   = t.EndsAt.Unix()
		period     = t.Period
		loop       = t.Loop
	)
	var err error

	if r.Form.Has("starts_at") {
		if startsUnix, err = strconv.ParseInt(r.Form.Get("starts_at"), 10, 0); err != nil {
			return middleware.HTTPError{
				Err:     err,
				Message: "Invalid starts_at",
				Code:    http.StatusBadRequest,
			}
		}
	}

	if r.Form.Has("ends_at") {
		if endsUnix, err = strconv.ParseInt(r.Form.Get("ends_at"), 10, 0); err != nil {
			return middleware.HTTPError{
				Err:     err,
				Message: "Invalid ends_at",
				Code:    http.StatusBadRequest,
			}
		}
	}

	if r.Form.Has("period_at") {
		if period, err = strconv.ParseInt(r.Form.Get("period"), 10, 0); err != nil {
			return middleware.HTTPError{
				Err:     err,
				Message: "Invalid period",
				Code:    http.StatusBadRequest,
			}
		}
	}

	if r.Form.Has("loop") {
		if loop, err = strconv.ParseInt(r.Form.Get("loop"), 10, 0); err != nil || loop <= 0 {
			return middleware.HTTPError{
				Err:     err,
				Message: "Invalid loop",
				Code:    http.StatusBadRequest,
			}
		}
	}

	if endsUnix < startsUnix {
		return middleware.HTTPError{
			Err:     err,
			Message: "End as earlier than start",
			Code:    http.StatusBadRequest,
		}
	}

	if endsUnix < 0 || startsUnix < 0 {
		return middleware.HTTPError{
			Err:     err,
			Message: "Timestamps cannot be less than 0",
			Code:    http.StatusBadRequest,
		}
	}

	if r.Form.Has("except") {
		exceptRaw := r.Form["except"]
		except := make([]int64, 0, len(exceptRaw))
		for _, s := range exceptRaw {
			if n, err := strconv.ParseInt(s, 10, 0); err == nil {
				except = append(except, n)
			}
		}
		t.Except = except
	}

	t.StartsAt = time.Unix(startsUnix, 0)
	t.EndsAt = time.Unix(endsUnix, 0)
	t.Period = period
	t.Loop = loop

	return nil
}
