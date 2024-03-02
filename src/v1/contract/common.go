package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetListParam struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Keyword string `json:"keyword"`
}

// ValidateQuery return common converted parameter from query parameter for get list data
// common query parameter is keyword, page, limit, and offset
// page is number page where the data is now, keyword is for search data by string keyword,
// limit is limit data loaded per page, offset is number data skiped when loaded data
// data page and limit from query parameter is always number in string
// its need to converted to int, it will return error if page and limit is not a number
func ValidateAndBuildRequest(r *http.Request) (getListParam *GetListParam, err error) {
	// default value for page and limit
	page, limit := 1, 10

	// get data from query parameter
	queryParams := r.URL.Query()
	limitQuery := queryParams.Get("limit")
	pageQuery := queryParams.Get("page")
	keyword := queryParams.Get("keyword")

	// query param validation
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			return
		}
	}

	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			return
		}
	}

	// offset for OFFSET in get list query
	offset := (page - 1) * limit
	getListParam = &GetListParam{
		Page:    page,
		Limit:   limit,
		Offset:  offset,
		Keyword: keyword,
	}

	return
}

func ValidateIDParamRequest(r *http.Request) (id int, err error) {
	idParam := chi.URLParam(r, "id")

	id, err = strconv.Atoi(idParam)
	if err != nil {
		log.Println(err)
		return id, err
	}

	return id, nil
}

func AddParameters(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func MarshalToReader(value interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonData), nil
}
