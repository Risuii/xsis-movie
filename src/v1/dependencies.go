package v1

import (
	"context"
	"log"

	"github.com/Risuii/movie/src/app"

	movieRepo "github.com/Risuii/movie/src/repository/movie"
	movieSvc "github.com/Risuii/movie/src/v1/service/movie"
)

type repositories struct {
	mRepo *movieRepo.MoviesRepository
}

type services struct {
	mSvc *movieSvc.MovieService
}

type Dependency struct {
	Repositories *repositories
	Services     *services
}

func initRepositories(ctx context.Context) *repositories {
	var r repositories
	var err error

	r.mRepo, err = movieRepo.InitMoviesRepository(ctx, app.DB(), app.Cache())
	if err != nil {
		log.Fatal("init movie repo err: ", err)
	}

	return &r
}

func initServices(ctx context.Context, r *repositories) *services {

	return &services{
		mSvc: movieSvc.InitMovieService(r.mRepo),
	}
}

func Dependencies(ctx context.Context) *Dependency {
	repositories := initRepositories(ctx)
	services := initServices(ctx, repositories)

	return &Dependency{
		Repositories: repositories,
		Services:     services,
	}
}
