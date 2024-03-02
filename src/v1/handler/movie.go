package handler

import (
	"log"
	"net/http"

	"github.com/Risuii/movie/src/errors"
	"github.com/Risuii/movie/src/middleware/response"
	"github.com/Risuii/movie/src/v1/contract"
)

func GetMovieHandler(svc MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := contract.ValidateIDParamRequest(r)
		if err != nil {
			log.Println(err)
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		data, err := svc.Get(r.Context(), id)
		if err != nil {
			log.Println(err)
			switch err {
			case errors.ErrMovieIdNotFound:
				response.JSONUnprocessableEntity(r.Context(), w, err)
			default:
				response.JSONInternalErrorResponse(r.Context(), w)
			}
			return
		}

		response.JSONSuccessResponse(r.Context(), w, data)
	}
}

func GetListMovieHandler(svc MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := contract.ValidateAndBuildRequest(r)
		if err != nil {
			log.Println(err)
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		data, err := svc.GetList(r.Context(), *params)
		if err != nil {
			log.Println(err)
			response.JSONInternalErrorResponse(r.Context(), w)
			return
		}

		response.JSONSuccessResponse(r.Context(), w, data)
	}
}

func CreateMovieHandler(svc MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movieRequest, err := contract.BuildAndValidateMovieRequest(r)
		if err != nil {
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		res, err := svc.Create(r.Context(), movieRequest)
		if err != nil {
			log.Println(err)
			response.JSONInternalErrorResponse(r.Context(), w)
			return
		}

		response.JSONSuccessResponse(r.Context(), w, res)
	}
}

func UpdateMovieHandler(svc MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := contract.ValidateIDParamRequest(r)
		if err != nil {
			log.Println(err)
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		movieRequest, err := contract.BuildAndValidateMovieRequest(r)
		if err != nil {
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		res, err := svc.Update(r.Context(), movieRequest, id)
		if err != nil {
			log.Println(err)
			switch err {
			case errors.ErrMovieIdNotFound:
				response.JSONUnprocessableEntity(r.Context(), w, err)
			default:
				response.JSONInternalErrorResponse(r.Context(), w)
			}
			return
		}

		response.JSONSuccessResponse(r.Context(), w, res)
	}
}

func DeleteMovieHandler(svc MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := contract.ValidateIDParamRequest(r)
		if err != nil {
			log.Println(err)
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			log.Println(err)
			switch err {
			case errors.ErrMovieIdNotFound:
				response.JSONUnprocessableEntity(r.Context(), w, err)
			default:
				response.JSONInternalErrorResponse(r.Context(), w)
			}
			return
		}

		response.JSONSuccessResponse(r.Context(), w, "success delete movie")
	}
}
