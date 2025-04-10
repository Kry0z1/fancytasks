package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/middleware"
	"github.com/Kry0z1/fancytasks/pkg/middleware/auth"
)

func LoginForToken(t auth.Tokenizer, h tasks.Hasher) func(http.ResponseWriter, *http.Request) error {
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
		_, err := auth.CheckUser(dctx, username, password, h)
		if err == auth.ErrInvalidCred || err == sql.ErrNoRows {
			return middleware.HTTPError{
				Err:     err,
				Message: "Wrong password or username",
				Code:    http.StatusUnauthorized,
			}
		}
		if err != nil {
			return err
		}

		token, err := t.CreateToken(map[string]any{
			"sub": username,
		}, 0)

		if err != nil {
			return err
		}

		w.Write([]byte(token))

		return nil
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) error {
	w.Write(tmpls["login.html"])
	return nil
}
