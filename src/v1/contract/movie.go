package contract

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	frsUtils "github.com/Risuii/frs-lib/utils"
	"github.com/go-playground/validator/v10"
)

type MovieResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating"`
	Image       string  `json:"image"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type MovieResponseDB struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Rating      float32   `db:"rating"`
	Image       string    `db:"image"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type GetListResponse struct {
	Data       []*MovieResponse
	Pagination *frsUtils.Pagination
}

type MovieRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating" validate:"required"`
	Image       string  `json:"image"`
}

func BuildAndValidateMovieRequest(r *http.Request) (MovieRequest, error) {
	var payload MovieRequest

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read request body err: ", err)
		return payload, err
	}

	if err := json.Unmarshal(bodyByte, &payload); err != nil {
		log.Println("unmarshal request body err: ", err)
		return payload, err
	}

	payload.Title = strings.ToLower(payload.Title)

	validator := validator.New()

	if err := validator.Struct(payload); err != nil {
		log.Println("validate request body err: ", err)
		return payload, err
	}

	return payload, nil
}
