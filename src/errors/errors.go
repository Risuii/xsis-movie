package errors

import (
	i18n_err "github.com/Risuii/frs-lib/i18n/errors"
)

var (
	ErrMovieIdNotFound = i18n_err.NewI18nError("err_movie_id_not_found")
	ErrDuplicatemovie  = i18n_err.NewI18nError("err_movie_duplicate")
)
