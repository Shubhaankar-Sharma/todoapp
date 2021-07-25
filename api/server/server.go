package server

import (
	"database/sql"
	"fmt"
	db "github.com/Shubhaankar-sharma/todoapp/db/sqlc"
	"github.com/Shubhaankar-sharma/todoapp/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Router *chi.Mux
	Queries *db.Queries
}


func NewServer(config utils.Config) *Server{
	s := &Server{}
	err := s.PrepareDB(config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	s.PrepareRouter()
	return s
}

func (s *Server) PrepareDB(config utils.Config) (err error){
	//Connect to db, else exit 0
	DB, err := sql.Open("postgres", config.DBUri)
	if err != nil {
		return
	}
	s.Queries = db.New(DB)
	log.Println("Connected to the Database")
	return
}

func (s *Server) PrepareRouter(){
	r := chi.NewRouter()

	//Use Global Middlewares Here
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	//Store Router in Struct
	s.Router = r
}

func (s *Server) RunServer(host string, port int) (err error) {

	log.Printf("Starting Server at %s:%v", host, port)

	err = http.ListenAndServe(fmt.Sprintf("%s:%v", host, port), s.Router)
	return
}
