package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	tasks "github.com/Kry0z1/fancytasks/pkg"
)

var contextUser struct{}

func getPopulatedContextWithUser(ctx context.Context, user *tasks.UserStored) context.Context {
	return context.WithValue(ctx, &contextUser, user)
}

func ContextUser(ctx context.Context) *tasks.UserStored {
	return ctx.Value(&contextUser).(*tasks.UserStored)
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
