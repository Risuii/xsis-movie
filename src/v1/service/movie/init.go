package movie

import (
	"context"

	"github.com/Risuii/movie/src/entity"
	"github.com/Risuii/movie/src/v1/contract"
)

type MovieRepository interface {
	Create(ctx context.Context, data *entity.Movie) (contract.MovieResponseDB, error)
	GetList(ctx context.Context, params contract.GetListParam) ([]*entity.Movie, error)
	GetMovieCount(ctx context.Context, param contract.GetListParam) (int64, error)
	Get(ctx context.Context, id int) (entity.Movie, error)
	Update(ctx context.Context, data *entity.Movie) error
	Delete(ctx context.Context, id int64) error
}
