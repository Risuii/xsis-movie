package v1

import (
	"net/http"

	"github.com/Risuii/movie/src/v1/handler"
	"github.com/go-chi/chi/v5"
)

func Router(r *chi.Mux, deps *Dependency) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Movie

	r.Route("/Movies", func(v1 chi.Router) {
		v1.Get("/{id}", handler.GetMovieHandler(deps.Services.mSvc))
		v1.Get("/", handler.GetListMovieHandler(deps.Services.mSvc))
		v1.Post("/", handler.CreateMovieHandler(deps.Services.mSvc))
		v1.Patch("/{id}", handler.UpdateMovieHandler(deps.Services.mSvc))
		v1.Delete("/{id}", handler.DeleteMovieHandler(deps.Services.mSvc))
	})
}
