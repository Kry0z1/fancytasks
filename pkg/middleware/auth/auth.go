package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/database"
)

var contextUser struct{}

func getPopulatedContextWithUser(ctx context.Context, user *tasks.User) context.Context {
	return context.WithValue(ctx, &contextUser, user)
}

func CheckUser(ctx context.Context, username, password string, hasher tasks.Hasher) (*tasks.User, error) {
	dctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))
	defer cancel()

	user, err := database.GetUserWithPassword(dctx, username)

	if err != nil {
		return nil, err
	}

	if !hasher.CheckPassword(password, user.HashedPassword) {
		return nil, ErrInvalidCred
	}

	return user, nil
}

func ContextUser(ctx context.Context) *tasks.User {
	return ctx.Value(&contextUser).(*tasks.User)
}

func CheckAuth(t Tokenizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", "Bearer")
			authorizationHeader := w.Header().Get("Authorization")
			if authorizationHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing Authorization header"))
				return
			}

			splitted := strings.Split(authorizationHeader, " ")
			if len(splitted) != 2 || splitted[0] != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid Authorization header format"))
				return
			}

			token := splitted[1]

			user, err := t.CheckToken(r.Context(), token)
			if err != nil {
				if errors.Is(err, ErrInvalidToken) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Invalid token"))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			next.ServeHTTP(
				w,
				r.WithContext(getPopulatedContextWithUser(r.Context(), user)),
			)
		})
	}
}
