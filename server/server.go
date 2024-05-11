package server

import (
	"net/http"
	"os"

	"context"
	apitoolkit "github.com/apitoolkit/apitoolkit-go"
	"github.com/apitoolkit/todo-backend-chi/server/controller"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// Server represents current server status
type Server struct {
	baseURL     string
	middlewares []func(w http.ResponseWriter, r *http.Request)
	router      *chi.Mux
}

// New creates a new Server with given URL
func New(baseURL string) *Server {
	return &Server{
		baseURL: baseURL,
		router:  chi.NewRouter(),
	}
}

func (s *Server) SetupRoutes(controller controller.Controller) {
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	apitoolkitClient, err := apitoolkit.NewClient(context.Background(), apitoolkit.Config{APIKey: os.Getenv("APITOOLKIT_KEY")})
	if err != nil {
		panic(err)
	}

	s.router.Use(apitoolkitClient.ChiMiddleware)
	s.router.Route("/", func(r chi.Router) {
		r.Get("/", controller.GetAll)   // Get all resources
		r.Post("/", controller.PostAll) // Create new resource

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", wrapperFunc(controller.GetOne))       // Get a single resource by ID
			r.Patch("/", wrapperFunc(controller.PatchOne))   // Update a resource by ID
			r.Post("/", wrapperFunc(controller.PostOne))     // Typically not used, added for completeness
			r.Delete("/", wrapperFunc(controller.DeleteOne)) // Delete a resource by ID
		})

		r.Options("/", controller.Options) // HTTP OPTIONS method
	})
}

func wrapperFunc(fn func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(ww http.ResponseWriter, rr *http.Request) {
		fn(ww, rr, chi.URLParam(rr, "id"))
	}
}

// Serve starts the actual serving job
func (s *Server) Serve(port string) {
	http.ListenAndServe(":"+port, s.router)
}
