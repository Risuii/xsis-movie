package movie

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	frsUtils "github.com/Risuii/frs-lib/utils"
	"github.com/Risuii/movie/src/app"
	"github.com/Risuii/movie/src/entity"
	"github.com/Risuii/movie/src/v1/contract"
	"github.com/go-faker/faker/v4"
	"github.com/mariomac/gostream/stream"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_movie "github.com/Risuii/movie/src/v1/service/mock/movie"
)

func TestMain(m *testing.M) {
	os.Chdir("../../../../")

	app.Init(context.Background())

	exitVal := m.Run()

	os.Exit(exitVal)

}

func TestGetMovieService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieRepo := mock_movie.NewMockMovieRepository(ctrl)

	type mockFields struct {
		movieRepo *mock_movie.MockMovieRepository
	}

	mocks := mockFields{
		movieRepo: mockMovieRepo,
	}

	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name     string
		args     args
		want     contract.MovieResponse
		wantErr  bool
		mockFunc func(mock mockFields, arg args)
	}{
		{
			name: "error starter pack id not found",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    contract.MovieResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.movieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, sql.ErrNoRows).Times(1)
			},
		},
		{
			name: "error starter pack error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    contract.MovieResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mock.movieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: contract.MovieResponse{
				CreatedAt: "0001-01-01 00:00:00",
				UpdatedAt: "0001-01-01 00:00:00",
			},
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mock.movieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			p := InitMovieService(mockMovieRepo)
			got, err := p.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Movie.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetListMovieService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieRepo := mock_movie.NewMockMovieRepository(ctrl)

	type mockFields struct {
		movieRepo *mock_movie.MockMovieRepository
	}

	mocks := mockFields{
		movieRepo: mockMovieRepo,
	}

	type args struct {
		ctx    context.Context
		params contract.GetListParam
	}

	var mockEntityMovie []*entity.Movie
	sizeDataset := 10
	for i := 0; i < sizeDataset; i++ {
		mockEntityMovie = append(mockEntityMovie, &entity.Movie{
			ModelID: entity.ModelID{
				Id: 1,
			},
			MovieData: entity.MovieData{
				Title:       faker.Name(),
				Description: faker.Paragraph(),
				Rating:      1,
				Image:       faker.Name(),
			},
		})
	}

	mockResponseMovieList := stream.Map(stream.OfSlice(mockEntityMovie), func(m *entity.Movie) *contract.MovieResponse {
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

	tests := []struct {
		name     string
		args     args
		want     contract.GetListResponse
		wantErr  bool
		mockFunc func(mock mockFields, args args)
	}{
		{
			name: "error get list",
			args: args{
				ctx: context.Background(),
				params: contract.GetListParam{
					Page:  1,
					Limit: 10,
				},
			},
			want:    contract.GetListResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, args args) {
				mockMovieRepo.EXPECT().GetList(gomock.Any(), args.params).Return([]*entity.Movie{}, assert.AnError).Times(1)
			},
		},
		{
			name: "error get count",
			args: args{
				ctx: context.Background(),
				params: contract.GetListParam{
					Page:  1,
					Limit: 10,
				},
			},
			want:    contract.GetListResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, args args) {
				mockMovieRepo.EXPECT().GetList(gomock.Any(), args.params).Return([]*entity.Movie{}, nil).Times(1)
				mockMovieRepo.EXPECT().GetMovieCount(gomock.Any(), args.params).Return(int64(1), assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				params: contract.GetListParam{
					Page:  1,
					Limit: 10,
				},
			},
			want: contract.GetListResponse{
				Data: mockResponseMovieList,
				Pagination: &frsUtils.Pagination{
					Page:      1,
					TotalPage: 1,
					TotalData: 1,
				},
			},
			wantErr: false,
			mockFunc: func(mock mockFields, args args) {
				mockMovieRepo.EXPECT().GetList(gomock.Any(), args.params).Return(mockEntityMovie, nil).Times(1)
				mockMovieRepo.EXPECT().GetMovieCount(gomock.Any(), args.params).Return(int64(1), nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			p := InitMovieService(mockMovieRepo)
			got, err := p.GetList(context.Background(), tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("movie.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateMovieService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieRepo := mock_movie.NewMockMovieRepository(ctrl)

	type mockFields struct {
		movieRepo *mock_movie.MockMovieRepository
	}

	mocks := mockFields{
		movieRepo: mockMovieRepo,
	}

	type args struct {
		ctx     context.Context
		request contract.MovieRequest
		params  *entity.Movie
	}

	tests := []struct {
		name     string
		args     args
		want     contract.MovieResponse
		wantErr  bool
		mockFunc func(mock mockFields, arg args)
	}{
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				request: contract.MovieRequest{},
				params:  &entity.Movie{},
			},
			want:    contract.MovieResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Create(gomock.Any(), arg.params).Return(contract.MovieResponseDB{}, assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				request: contract.MovieRequest{
					Title:       "test-title-1",
					Description: "test-description-1",
					Rating:      1,
					Image:       "test-image-1",
				},
				params: &entity.Movie{
					MovieData: entity.MovieData{
						Title:       "test-title-1",
						Description: "test-description-1",
						Rating:      1,
						Image:       "test-image-1",
					},
				},
			},
			want: contract.MovieResponse{
				CreatedAt: "0001-01-01 00:00:00",
				UpdatedAt: "0001-01-01 00:00:00",
			},
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Create(gomock.Any(), arg.params).Return(contract.MovieResponseDB{}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			p := InitMovieService(mockMovieRepo)
			got, err := p.Create(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Movie.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateMovieService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieRepo := mock_movie.NewMockMovieRepository(ctrl)

	type mockFields struct {
		movieRepo *mock_movie.MockMovieRepository
	}

	mocks := mockFields{
		movieRepo: mockMovieRepo,
	}

	type args struct {
		ctx     context.Context
		request contract.MovieRequest
		params  *entity.Movie
		id      int
	}

	tests := []struct {
		name     string
		args     args
		want     contract.MovieResponse
		wantErr  bool
		mockFunc func(mock mockFields, arg args)
	}{
		{
			name: "error id not found",
			args: args{
				ctx:     context.Background(),
				request: contract.MovieRequest{},
				params:  &entity.Movie{},
				id:      1,
			},
			want:    contract.MovieResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, sql.ErrNoRows).Times(1)
			},
		},
		{
			name: "error update",
			args: args{
				ctx:     context.Background(),
				request: contract.MovieRequest{},
				params:  &entity.Movie{},
				id:      1,
			},
			want:    contract.MovieResponse{},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, nil).Times(1)
				mockMovieRepo.EXPECT().Update(gomock.Any(), arg.params).Return(assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				request: contract.MovieRequest{},
				params:  &entity.Movie{},
				id:      1,
			},
			want: contract.MovieResponse{
				CreatedAt: "0001-01-01 00:00:00",
				UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, nil).Times(1)
				mockMovieRepo.EXPECT().Update(gomock.Any(), arg.params).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			p := InitMovieService(mockMovieRepo)
			got, err := p.Update(tt.args.ctx, tt.args.request, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Movie.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteMovieService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieRepo := mock_movie.NewMockMovieRepository(ctrl)

	type mockFields struct {
		movieRepo *mock_movie.MockMovieRepository
	}

	mocks := mockFields{
		movieRepo: mockMovieRepo,
	}

	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		mockFunc func(mock mockFields, arg args)
	}{
		{
			name: "error id not found",
			args: args{
				id: 1,
			},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, sql.ErrNoRows).Times(1)
			},
		},
		{
			name: "error delete",
			args: args{
				id: 1,
			},
			wantErr: true,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, nil).Times(1)
				mockMovieRepo.EXPECT().Delete(gomock.Any(), int64(0)).Return(assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				id: 1,
			},
			wantErr: false,
			mockFunc: func(mock mockFields, arg args) {
				mockMovieRepo.EXPECT().Get(gomock.Any(), arg.id).Return(entity.Movie{}, nil).Times(1)
				mockMovieRepo.EXPECT().Delete(gomock.Any(), int64(0)).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)

			p := InitMovieService(mockMovieRepo)
			err := p.Delete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Movie.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
