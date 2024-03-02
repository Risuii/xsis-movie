package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Risuii/movie/src/v1/contract"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	appErr "github.com/Risuii/movie/src/errors"
	mock_handler "github.com/Risuii/movie/src/v1/handler/mock"
)

func TestGetMovieHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieSvc := mock_handler.NewMockMovieService(ctrl)

	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name       string
		args       args
		mockFunc   func(arg args)
		want       contract.MovieResponse
		wantErr    bool
		statusCode int
		parameter  map[string]string
	}{
		{
			name: "error bad request",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
			parameter:  nil,
			mockFunc:   func(arg args) {},
		},
		{
			name: "error id not found",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusUnprocessableEntity,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Get(gomock.Any(), arg.id).Return(contract.MovieResponse{}, appErr.ErrMovieIdNotFound).Times(1)
			},
		},
		{
			name: "error internal server",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Get(gomock.Any(), arg.id).Return(contract.MovieResponse{}, assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       contract.MovieResponse{},
			wantErr:    false,
			statusCode: http.StatusOK,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Get(gomock.Any(), arg.id).Return(contract.MovieResponse{}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)

			req, err := http.NewRequest(http.MethodGet, "/just/for/testing", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = contract.AddParameters(req, tt.parameter)

			r := httptest.NewRecorder()
			handler := http.HandlerFunc(GetMovieHandler(mockMovieSvc))
			handler.ServeHTTP(r, req)

			if r.Code != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", r.Code, tt.statusCode)
			}
		})
	}
}

func TestGetListMovieHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieSvc := mock_handler.NewMockMovieService(ctrl)

	mockParams := contract.GetListParam{
		Page:   1,
		Limit:  10,
		Offset: 0,
	}

	type args struct {
		ctx    context.Context
		params contract.GetListParam
	}

	tests := []struct {
		name       string
		args       args
		mockFunc   func(arg args)
		want       contract.GetListResponse
		wantErr    bool
		statusCode int
	}{
		{
			name: "error",
			args: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			want:       contract.GetListResponse{},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().GetList(gomock.Any(), mockParams).Return(contract.GetListResponse{}, assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			want:       contract.GetListResponse{},
			wantErr:    false,
			statusCode: http.StatusOK,
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().GetList(gomock.Any(), mockParams).Return(contract.GetListResponse{}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)

			req, err := http.NewRequest(http.MethodGet, "/just/for/testing", nil)
			if err != nil {
				t.Fatal(err)
			}

			r := httptest.NewRecorder()
			handler := http.HandlerFunc(GetListMovieHandler(mockMovieSvc))
			handler.ServeHTTP(r, req)

			if r.Code != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", r.Code, tt.statusCode)
			}
		})
	}
}

func TestCreateMovieHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieSvc := mock_handler.NewMockMovieService(ctrl)

	mockRequest := contract.MovieRequest{
		Title:       "test-name",
		Description: "test-description",
		Rating:      1,
		Image:       "test-image",
	}

	type args struct {
		ctx    context.Context
		params contract.MovieRequest
	}

	tests := []struct {
		name       string
		args       args
		mockFunc   func(arg args)
		want       contract.MovieResponse
		wantErr    bool
		statusCode int
	}{
		{
			name: "error bad request",
			args: args{
				ctx: context.Background(),
				params: contract.MovieRequest{
					Title:       "",
					Description: "",
					Rating:      0,
					Image:       "",
				},
			},
			mockFunc:   func(arg args) {},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "error internal server",
			args: args{
				ctx:    context.Background(),
				params: mockRequest,
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Create(gomock.Any(), arg.params).Return(contract.MovieResponse{}, assert.AnError).Times(1)
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				params: mockRequest,
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Create(gomock.Any(), arg.params).Return(contract.MovieResponse{}, nil).Times(1)
			},
			want:       contract.MovieResponse{},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)

			reader, err := contract.MarshalToReader(tt.args.params)
			if err != nil {
				t.Errorf("Error when try to marshal params. error = %v, data = %v", err, tt.args.params)
				return
			}
			req, err := http.NewRequest(http.MethodPost, "/just/for/testing", reader)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateMovieHandler(mockMovieSvc))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.statusCode)
			}
		})
	}
}

func TestUpdateMovieHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieSvc := mock_handler.NewMockMovieService(ctrl)

	mockRequest := contract.MovieRequest{
		Title:       "test-name",
		Description: "test-description",
		Rating:      1,
		Image:       "test-image",
	}

	type args struct {
		ctx    context.Context
		id     int
		params contract.MovieRequest
	}

	tests := []struct {
		name       string
		args       args
		mockFunc   func(arg args)
		want       contract.MovieResponse
		wantErr    bool
		statusCode int
		parameter  map[string]string
	}{
		{
			name: "error bad request id",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
			parameter:  nil,
			mockFunc:   func(arg args) {},
		},
		{
			name: "error bad request payload",
			args: args{
				ctx: context.Background(),
				id:  1,
				params: contract.MovieRequest{
					Title:       "",
					Description: "",
					Rating:      0,
					Image:       "",
				},
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {},
		},
		{
			name: "error internal server",
			args: args{
				ctx:    context.Background(),
				id:     1,
				params: mockRequest,
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Update(gomock.Any(), arg.params, arg.id).Return(contract.MovieResponse{}, assert.AnError).Times(1)
			},
		},
		{
			name: "error id not found",
			args: args{
				ctx:    context.Background(),
				id:     1,
				params: mockRequest,
			},
			want:       contract.MovieResponse{},
			wantErr:    true,
			statusCode: http.StatusUnprocessableEntity,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Update(gomock.Any(), arg.params, arg.id).Return(contract.MovieResponse{}, appErr.ErrMovieIdNotFound).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				id:     1,
				params: mockRequest,
			},
			want:       contract.MovieResponse{},
			wantErr:    false,
			statusCode: http.StatusOK,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Update(gomock.Any(), arg.params, arg.id).Return(contract.MovieResponse{}, nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)

			reader, err := contract.MarshalToReader(tt.args.params)
			if err != nil {
				t.Errorf("Error when try to marshal params. error = %v, data = %v", err, tt.args.params)
				return
			}
			req, err := http.NewRequest(http.MethodPost, "/just/for/testing", reader)
			if err != nil {
				t.Fatal(err)
			}

			req = contract.AddParameters(req, tt.parameter)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(UpdateMovieHandler(mockMovieSvc))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.statusCode)
			}
		})
	}
}

func TestDeleteMovieHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMovieSvc := mock_handler.NewMockMovieService(ctrl)

	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name       string
		args       args
		mockFunc   func(arg args)
		want       string
		wantErr    bool
		statusCode int
		parameter  map[string]string
	}{
		{
			name: "error bad request",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       "",
			wantErr:    true,
			statusCode: http.StatusBadRequest,
			parameter:  nil,
			mockFunc:   func(arg args) {},
		},
		{
			name: "error id not found",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       "",
			wantErr:    true,
			statusCode: http.StatusUnprocessableEntity,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Delete(gomock.Any(), arg.id).Return(appErr.ErrMovieIdNotFound).Times(1)
			},
		},
		{
			name: "error internal server",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       "",
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Delete(gomock.Any(), arg.id).Return(assert.AnError).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:       "",
			wantErr:    true,
			statusCode: http.StatusOK,
			parameter: map[string]string{
				"id": "1",
			},
			mockFunc: func(arg args) {
				mockMovieSvc.EXPECT().Delete(gomock.Any(), arg.id).Return(nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.args)

			req, err := http.NewRequest(http.MethodGet, "/just/for/testing", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = contract.AddParameters(req, tt.parameter)

			r := httptest.NewRecorder()
			handler := http.HandlerFunc(DeleteMovieHandler(mockMovieSvc))
			handler.ServeHTTP(r, req)

			if r.Code != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", r.Code, tt.statusCode)
			}
		})
	}
}
