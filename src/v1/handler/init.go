package handler

import (
	"context"

	"github.com/Risuii/movie/src/v1/contract"
)

type MovieService interface {
	Get(ctx context.Context, id int) (res contract.MovieResponse, err error)
	GetList(ctx context.Context, params contract.GetListParam) (res contract.GetListResponse, err error)
	Create(ctx context.Context, request contract.MovieRequest) (res contract.MovieResponse, err error)
	Update(ctx context.Context, request contract.MovieRequest, id int) (res contract.MovieResponse, err error)
	Delete(ctx context.Context, id int) (err error)
}
