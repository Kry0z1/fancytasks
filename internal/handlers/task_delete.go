package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/database"
	"github.com/Kry0z1/fancytasks/pkg/middleware"
	"github.com/Kry0z1/fancytasks/pkg/middleware/auth"
)

func DeleteTask(w http.ResponseWriter, r *http.Request) error {
	var (
		id       int
		err      error
		delete   func() error
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
		delete = func() error { return database.DeleteBaseTask(dctx, &task) }
		baseTask = &task
	case "event":
		var task tasks.Event
		delete = func() error { return database.DeleteEvent(dctx, &task) }
		baseTask = &task.BaseTask
	case "deadline":
		var task tasks.TaskWithDeadline
		delete = func() error { return database.DeleteTaskWithDeadline(dctx, &task) }
		baseTask = &task.BaseTask
	case "repeat":
		var task tasks.RepeatingTask
		delete = func() error { return database.DeleteRepeatingTask(dctx, &task) }
		baseTask = &task.BaseTask
	default:
		return middleware.HTTPError{
			Err:     nil,
			Message: "Invalid task type",
			Code:    http.StatusBadRequest,
		}
	}

	baseTask.ID = id
	err = delete()
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

	w.Write([]byte("Successful"))
	return nil
}
