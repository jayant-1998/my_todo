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
			utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
				"status": "server is running",
			})
		})
		v1.Route("/public", func(p chi.Router) {
			p.Post("/register", handler.CreateUser)
			p.Post("/login", handler.LoginUser)
		})
		v1.Route("/user", func(u chi.Router) {
			u.Use(middlewares.AuthMiddleware)
			u.Get("/", handler.UserInfo)
			u.Delete("/", handler.DeleteUser)
			u.Delete("/logout", handler.Logout)
			u.Put("/", handler.UpdateUser)
		})

		v1.Route("/task", func(t chi.Router) {
			t.Use(middlewares.AuthMiddleware)
			t.Post("/", handler.CreateTask)
			t.Get("/", handler.AllTask)
			t.Delete("/{taskId}", handler.DeleteTask)
			t.Put("/", handler.UpdateTask)
			t.Put("/{taskId}/complete", handler.Complete)
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
