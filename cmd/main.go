package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kry0z1/fancytasks/internal/handlers"
	_ "github.com/Kry0z1/fancytasks/internal/redis"
	"github.com/Kry0z1/fancytasks/pkg/middleware"
)

func main() {
	http.Handle("GET /", middleware.CollectErrorFunc(handlers.Index, middleware.Logger))
	http.Handle("POST /", middleware.Logger(http.RedirectHandler(":8000/register", http.StatusFound)))
	http.Handle("GET /register", middleware.CollectErrorFunc(handlers.RegisterPage, middleware.Logger))
	http.Handle("POST /register", middleware.CollectErrorFunc(handlers.Register, middleware.Logger))

	fmt.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
