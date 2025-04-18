package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Kry0z1/fancytasks/internal/middleware"
	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/database"
)

func Register(h tasks.Hasher) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}

		usernameSlice := r.Form["username"]
		passwordSlice := r.Form["password"]

		if len(usernameSlice) == 0 || len(passwordSlice) == 0 {
			return middleware.HTTPError{
				Err:     nil,
				Message: "Missing password or username",
				Code:    http.StatusNotFound,
			}
		}

		username := usernameSlice[0]
		password := passwordSlice[0]

		if username == "" || password == "" {
			return middleware.HTTPError{
				Err:     nil,
				Message: "Missing password or username",
				Code:    http.StatusNotFound,
			}
		}

		dctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(time.Second))
		defer cancel()
		_, err := database.CreateUser(dctx, username, password, h)

		if err == database.ErrUserExists {
			return middleware.HTTPError{
				Err:     database.ErrUserExists,
				Message: "Invalid password or username",
				Code:    http.StatusUnauthorized,
			}
		}

		if err != nil {
			return err
		}

		w.Write([]byte("ok"))

		return nil
	}
}
