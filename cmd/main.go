package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kry0z1/fancytasks/internal/handlers"
	_ "github.com/Kry0z1/fancytasks/internal/redis"
	tasks "github.com/Kry0z1/fancytasks/pkg"
	"github.com/Kry0z1/fancytasks/pkg/middleware"
	"github.com/Kry0z1/fancytasks/pkg/middleware/auth"
)

func main() {
	t, err := auth.NewTokenizer(tasks.Cfg.JWT.GetExpiresDelta(), os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Couldn't create tokenizer: %s", err.Error())
	}
	h := tasks.NewHasher()

	http.Handle("GET /", middleware.CollectErrorFunc(handlers.Index, middleware.Logger))
	http.Handle("GET /register", middleware.CollectErrorFunc(handlers.RegisterPage, middleware.Logger))
	http.Handle("POST /register", middleware.CollectErrorFunc(handlers.Register(h), middleware.Logger))
	http.Handle("GET /login", middleware.CollectErrorFunc(handlers.LoginPage, middleware.Logger))
	http.Handle("POST /login", middleware.CollectErrorFunc(handlers.LoginForToken(t, h), middleware.Logger))
	http.Handle("GET /secret", middleware.CollectErrorFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("ok"))
		return nil
	}, middleware.Logger, auth.CheckAuth(t)))

	fmt.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
