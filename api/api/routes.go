package api

import (
	"github.com/Shubhaankar-sharma/todoapp/api/handlers"
	db "github.com/Shubhaankar-sharma/todoapp/db/sqlc"
	"github.com/Shubhaankar-sharma/todoapp/utils"
	"github.com/go-chi/chi/v5"
)

func MainRouter(r *chi.Mux, queries *db.Queries, config utils.Config) {
	r.Use(SetContentTypeMiddleware)
	r.Group(func(r chi.Router) {
		r.Get("/", handlers.HelloWorld())
		r.Post("/auth/register", handlers.CreateUserHandler(queries, config.SecretKey))
		r.Post("/auth/login", handlers.Login(queries, config.SecretKey))

	})

	//Private Routes
	r.Group(func(r chi.Router) {
		r.Use(AuthJwtWrap(config.SecretKey))
		r.Get("/testing", handlers.HelloWorld())
		r.Post("/todo/create", handlers.CreateToDoHandler(queries))
	})
}