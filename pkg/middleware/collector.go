package middleware

import "net/http"

func CollectFunc(
	f func(http.ResponseWriter, *http.Request),
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	var result http.Handler = http.HandlerFunc(f)
	for _, m := range middlewares {
		result = m(result)
	}

	return result
}

func CollectErrorFunc(
	f func(http.ResponseWriter, *http.Request) error,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	result := ErrorMiddleware(f)
	for _, m := range middlewares {
		result = m(result)
	}

	return result
}

func Collect(
	f http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}
