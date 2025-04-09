package middleware

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Err     error
	Message string
	Code    int
}

func (he HTTPError) Error() string {
	return fmt.Errorf("%w: %s", he.Err, he.Message).Error()
}

func ErrorMiddleware(next func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)

		if err == nil {
			return
		}

		var he *HTTPError
		if errors.As(err, &he) {
			http.Error(w, he.Message, he.Code)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})
}
