package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Risuii/movie/src/app"
	"github.com/Risuii/movie/src/middleware/request"
	"github.com/go-chi/chi/v5"

	v1 "github.com/Risuii/movie/src/v1"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	initCtx := context.Background()
	if err := app.Init(initCtx); err != nil {
		panic(err)
	}

	startService(initCtx)
}

func startService(ctx context.Context) {
	address := fmt.Sprintf(":%d", app.Config().BindAddress)

	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(request.RequestIDContext(request.DefaultGenerator))
	r.Use(request.RequestAttributesContext)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	deps := v1.Dependencies(ctx)
	v1.Router(r, deps)

	err := http.ListenAndServe(address, r)
	if err != nil {
		log.Println(err)
	}
}
