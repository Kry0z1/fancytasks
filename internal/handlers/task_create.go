package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Kry0z1/fancytasks/internal/middleware"
	"github.com/Kry0z1/fancytasks/internal/middleware/auth"
	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/database"
)

func CreateTask(w http.ResponseWriter, r *http.Request) error {
	user := auth.ContextUser(r.Context())
	if user == nil {
		return middleware.HTTPError{
			Err:     nil,
			Message: "Unauthorized",
			Code:    http.StatusUnauthorized,
		}
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	taskType := r.Form.Get("tasktype")
	var endFunc func(http.ResponseWriter, *http.Request, *tasks.BaseTask) error

	switch taskType {
	case "":
		return middleware.HTTPError{
			Err:     nil,
			Message: "Task type not found in form",
			Code:    http.StatusNotFound,
		}
	case "basetask":
		endFunc = func(w http.ResponseWriter, r *http.Request, t *tasks.BaseTask) error {
			if err := database.CreateBaseTask(r.Context(), t); err != nil {
				return err
			}
			return json.NewEncoder(w).Encode(t)
		}
	case "event":
		endFunc = createEvent
	case "deadline":
		endFunc = createDeadline
	case "repeat":
		endFunc = createRepeatingTask
	default:
		return middleware.HTTPError{
			Err:     nil,
			Message: "Invalid task type",
			Code:    http.StatusBadRequest,
		}
	}

	task, err := createBaseTask(r)
	if err != nil {
		return err
	}
	return endFunc(w, r, task)
}

func createBaseTask(r *http.Request) (*tasks.BaseTask, error) {
	var task tasks.BaseTask

	if task.Title = r.Form.Get("title"); task.Title == "" {
		return nil, middleware.HTTPError{
			Err:     nil,
			Message: "Missing task title",
			Code:    http.StatusNotFound,
		}
	}

	task.Description = r.Form.Get("description")
	task.Owner = auth.ContextUser(r.Context()).Username
	task.Topic = r.Form.Get("topic")
	if task.Topic == "" {
		task.Topic = "default"
	}

	return &task, nil
}

func createEvent(w http.ResponseWriter, r *http.Request, t *tasks.BaseTask) error {
	var startsUnix, endsUnix int64
	var err error

	if startsUnix, err = strconv.ParseInt(r.Form.Get("starts_at"), 10, 0); err != nil {
		return middleware.HTTPError{
			Err:     err,
			Message: "Invalid starts_at",
			Code:    http.StatusBadRequest,
		}
	}

	if endsUnix, err = strconv.ParseInt(r.Form.Get("ends_at"), 10, 0); err != nil {
		return middleware.HTTPError{
			Err:     err,
			Message: "Invalid ends_at",
			Code:    http.StatusBadRequest,
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

	result := tasks.Event{
		BaseTask: *t,
		StartsAt: time.Unix(startsUnix, 0),
		EndsAt:   time.Unix(endsUnix, 0),
	}

	if err := database.CreateEvent(r.Context(), &result); err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(result)
}

func createDeadline(w http.ResponseWriter, r *http.Request, t *tasks.BaseTask) error {
	var deadline int64
	var err error

	if deadline, err = strconv.ParseInt(r.Form.Get("deadline"), 10, 0); err != nil {
		return middleware.HTTPError{
			Err:     err,
			Message: "Invalid deadline",
			Code:    http.StatusBadRequest,
		}
	}

	if deadline < 0 {
		return middleware.HTTPError{
			Err:     err,
			Message: "Timestamps cannot be less than 0",
			Code:    http.StatusBadRequest,
		}
	}

	result := tasks.TaskWithDeadline{
		BaseTask: *t,
		Deadline: time.Unix(deadline, 0),
	}

	if err := database.CreateTaskWithDeadline(r.Context(), &result); err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(result)
}

func createRepeatingTask(w http.ResponseWriter, r *http.Request, t *tasks.BaseTask) error {
	var startsUnix, endsUnix, period, loop int64
	var err error

	if startsUnix, err = strconv.ParseInt(r.Form.Get("starts_at"), 10, 0); err != nil {
		return middleware.HTTPError{
			Err:     err,
			Message: "Invalid starts_at",
			Code:    http.StatusBadRequest,
		}
	}

	if endsUnix, err = strconv.ParseInt(r.Form.Get("ends_at"), 10, 0); err != nil {
		return middleware.HTTPError{
			Err:     err,
			Message: "Invalid ends_at",
			Code:    http.StatusBadRequest,
		}
	}

	if period, err = strconv.ParseInt(r.Form.Get("period"), 10, 0); err != nil {
		return middleware.HTTPError{
			Err:     err,
			Message: "Invalid period",
			Code:    http.StatusBadRequest,
		}
	}

	if loop, err = strconv.ParseInt(r.Form.Get("loop"), 10, 0); err != nil || loop <= 0 {
		return middleware.HTTPError{
			Err:     err,
			Message: "Invalid loop",
			Code:    http.StatusBadRequest,
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

	exceptRaw := r.Form["except"]
	except := make([]int64, 0, len(exceptRaw))
	for _, s := range exceptRaw {
		if n, err := strconv.ParseInt(s, 10, 0); err == nil {
			except = append(except, n)
		}
	}

	result := tasks.RepeatingTask{
		Event: tasks.Event{
			BaseTask: *t,
			StartsAt: time.Unix(startsUnix, 0),
			EndsAt:   time.Unix(endsUnix, 0),
		},
		Period: period,
		Loop:   loop,
		Except: except,
	}

	if err := database.CreateRepeatingTask(r.Context(), &result); err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(result)
}
