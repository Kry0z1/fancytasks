package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Kry0z1/fancytasks/pkg/database"
	"github.com/Kry0z1/fancytasks/pkg/middleware"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) error {
	w.Write(tmpls["register.html"])
	return nil
}

func Register(w http.ResponseWriter, r *http.Request) error {
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

	dctx, _ := context.WithDeadline(r.Context(), time.Now().Add(time.Second))
	_, err := database.CreateUser(dctx, username, password)

	if err == database.ErrUserExists {
		return middleware.HTTPError{
			Err:     nil,
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
