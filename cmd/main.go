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

// TODO:
// Добавить удаление тасок
// Добавить разделение тасок на топики
// Мигрировать дб (добавить поле isAdmin)
// Сделать кэширование запросов с помощью redis

func main() {
	t, err := auth.NewTokenizer(tasks.Cfg.JWT.GetExpiresDelta(), os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Couldn't create tokenizer: %s", err.Error())
	}
	h := tasks.NewHasher()

	http.Handle("POST /register", middleware.LoggerErrorFunc(handlers.Register(h)))
	http.Handle("POST /login", middleware.LoggerErrorFunc(handlers.LoginForToken(t, h)))
	http.Handle("GET /tasks", middleware.LoggerAuthErrorFunc(handlers.Me, t))
	http.Handle("POST /tasks/create", middleware.LoggerAuthErrorFunc(handlers.CreateTask, t))
	http.Handle("PUT /tasks/update", middleware.LoggerAuthErrorFunc(handlers.UpdateTask, t))
	http.Handle("DELETE /tasks/delete", middleware.LoggerAuthErrorFunc(handlers.DeleteTask, t))
	http.Handle("GET /secret", middleware.LoggerAuthErrorFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("ok"))
		return nil
	}, t))

	fmt.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
