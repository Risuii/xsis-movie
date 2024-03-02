package movie

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Risuii/movie/src/entity"
	"github.com/Risuii/movie/src/v1/contract"
	"github.com/mariomac/gostream/stream"

	frsUtils "github.com/Risuii/frs-lib/utils"
	appErr "github.com/Risuii/movie/src/errors"
)

type MovieService struct {
	MovieRepo MovieRepository
}

func InitMovieService(mRepo MovieRepository) *MovieService {
	return &MovieService{
		MovieRepo: mRepo,
	}
}

func useNewValueIfNotNull(newValue, oldValue string) string {
	if newValue != "" {
		return newValue
	}
	return oldValue
}

func useNewFloatValueIfNotZero(newValue, oldValue float32) float32 {
	if newValue != 0 {
		return newValue
	}
	return oldValue
}

func mapperMovieRequest(movie *entity.Movie, request *contract.MovieRequest) *entity.Movie {
	movie.Title = useNewValueIfNotNull(request.Title, movie.Title)
	movie.Description = useNewValueIfNotNull(request.Description, movie.Description)
	movie.Rating = useNewFloatValueIfNotZero(request.Rating, movie.Rating)
	movie.Image = useNewValueIfNotNull(request.Image, movie.Image)

	return movie
}

func (ms *MovieService) Get(ctx context.Context, id int) (res contract.MovieResponse, err error) {

	movie, err := ms.MovieRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = appErr.ErrMovieIdNotFound
		}
		log.Println("get movie err: ", err)
		return
	}

	res = contract.MovieResponse{
		ID:          int(movie.Id),
		Title:       movie.Title,
		Description: movie.Description,
		Rating:      movie.Rating,
		Image:       movie.Image,
		CreatedAt:   movie.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   movie.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return
}

func (ms *MovieService) GetList(ctx context.Context, params contract.GetListParam) (res contract.GetListResponse, err error) {

	movie, err := ms.MovieRepo.GetList(ctx, params)
	if err != nil {
		log.Println("get list movie err: ", err)
		return
	}

	count, err := ms.MovieRepo.GetMovieCount(ctx, params)
	if err != nil {
		log.Println("get count movie err: ", err)
		return
	}

	pagination := frsUtils.GetPaginationData(params.Page, params.Limit, int(count))

	responseMovieList := stream.Map(stream.OfSlice(movie), func(m *entity.Movie) *contract.MovieResponse {
		return &contract.MovieResponse{
			ID:          int(m.Id),
			Title:       m.Title,
			Description: m.Description,
			Rating:      m.Rating,
			Image:       m.Image,
			CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   m.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}).ToSlice()

	res = contract.GetListResponse{
		Data:       responseMovieList,
		Pagination: pagination,
	}

	return
}

func (ms *MovieService) Create(ctx context.Context, request contract.MovieRequest) (res contract.MovieResponse, err error) {

	req := &entity.Movie{
		MovieData: entity.MovieData{
			Title:       request.Title,
			Description: request.Description,
			Rating:      request.Rating,
			Image:       request.Image,
		},
	}

	movie, err := ms.MovieRepo.Create(ctx, req)
	if err != nil {
		log.Println("error create movie err: ", err)
		return
	}

	res = contract.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Rating:      movie.Rating,
		Image:       movie.Image,
		CreatedAt:   movie.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   movie.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return
}

func (ms *MovieService) Update(ctx context.Context, request contract.MovieRequest, id int) (res contract.MovieResponse, err error) {

	movie, err := ms.MovieRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = appErr.ErrMovieIdNotFound
		}
		log.Println("find movie err: ", err)
		return
	}

	movie = *mapperMovieRequest(&movie, &request)

	err = ms.MovieRepo.Update(ctx, &movie)
	if err != nil {
		log.Println("update movie err: ", err)
		return
	}

	res = contract.MovieResponse{
		ID:          int(movie.Id),
		Title:       movie.Title,
		Description: movie.Description,
		Rating:      movie.Rating,
		Image:       movie.Image,
		CreatedAt:   movie.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	return
}

func (ms *MovieService) Delete(ctx context.Context, id int) (err error) {

	movie, err := ms.MovieRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = appErr.ErrMovieIdNotFound
		}
		log.Println("get movie err: ", err)
		return
	}

	err = ms.MovieRepo.Delete(ctx, movie.Id)
	if err != nil {
		log.Println("delete err: ", err)
		return
	}

	return
}
