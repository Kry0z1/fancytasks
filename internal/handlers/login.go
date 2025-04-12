package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
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

		username := r.FormValue("username")
		password := r.FormValue("password")

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

		return json.NewEncoder(w).Encode(map[string]string{
			"token_type":   "Bearer",
			"access_token": token,
		})
	}
}
