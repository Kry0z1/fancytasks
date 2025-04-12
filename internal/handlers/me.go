package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"time"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/database"
	"github.com/Kry0z1/fancytasks/pkg/middleware"
	"github.com/Kry0z1/fancytasks/pkg/middleware/auth"
)

func Me(w http.ResponseWriter, r *http.Request) error {
	user := auth.ContextUser(r.Context())
	if user == nil {
		return middleware.HTTPError{
			Err:     nil,
			Message: "Unauthorized",
			Code:    http.StatusUnauthorized,
		}
	}

	filter := r.URL.Query()["filter"]
	var userDB *tasks.User = &tasks.User{Username: user.Username}
	var err error

	dctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(time.Second))
	defer cancel()

	var found bool

	if slices.Contains(filter, "base") {
		userDB.BaseTasks, err = database.GetUserBaseTasks(dctx, user.Username)
		found = true
	}
	if slices.Contains(filter, "events") {
		userDB.Events, err = database.GetUserEvents(dctx, user.Username)
		found = true
	}
	if slices.Contains(filter, "repeat") {
		userDB.RepeatingTasks, err = database.GetUserRepeatingTasks(dctx, user.Username)
		found = true
	}
	if slices.Contains(filter, "deadline") {
		userDB.TasksWithDeadline, err = database.GetUserTasksWithDeadline(dctx, user.Username)
		found = true
	}

	if !found {
		userDB, err = database.GetUserWithTasks(dctx, user.Username)
	}

	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(userDB)
}
