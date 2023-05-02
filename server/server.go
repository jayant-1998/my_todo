package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"jayant/handler"
	"jayant/middlewares"
	"jayant/utils"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRouter() *Server {
	router := chi.NewRouter()
	router.Route("/todo", func(v1 chi.Router) {
		v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "server is running"})
		})
		v1.Route("/", func(c chi.Router) {
			c.Post("/register", handler.CreateUser)
			c.Post("/login", handler.LoginUser)
		})
		v1.Route("/user", func(public chi.Router) {
			public.Use(middlewares.AuthMiddleware)
			public.Get("/", handler.InfoUser)
			public.Delete("/", handler.DeleteUser)
			public.Delete("/logout", handler.Logout)
			public.Put("/", handler.UpdateUser)
		})

		v1.Route("/task", func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware)
			r.Post("/", handler.CreateTask)
			r.Get("/", handler.GetTodoInfo)
			r.Delete("/{name}", handler.DeleteTodo)
			r.Put("/", handler.UpdateTodo)
			r.Put("/complete/{task}", handler.IsComplete)
		})
	})
	return &Server{
		Router: router,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
