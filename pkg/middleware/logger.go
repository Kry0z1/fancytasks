package middleware

import (
	"log"
	"net/http"
)

type ResponseWriterWithStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (e *ResponseWriterWithStatusCode) WriteHeader(statusCode int) {
	e.statusCode = statusCode
	e.ResponseWriter.WriteHeader(statusCode)
}

// If status code is not set yet, returns 200
func (e ResponseWriterWithStatusCode) StatusCode() int {
	if e.statusCode == 0 {
		return http.StatusOK
	}
	return e.statusCode
}

func ExtendResponseWriter(w http.ResponseWriter) *ResponseWriterWithStatusCode {
	return &ResponseWriterWithStatusCode{w, 0}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ew := ExtendResponseWriter(w)

		next.ServeHTTP(ew, r)

		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
		log.Printf("%s %d %s\n", r.Method, ew.StatusCode(), r.URL)
	})
}
