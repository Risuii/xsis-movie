package movie

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Risuii/movie/src/entity"
	"github.com/Risuii/movie/src/v1/contract"
)

func (mr *MoviesRepository) GetList(ctx context.Context, params contract.GetListParam) ([]*entity.Movie, error) {
	var Movie []*entity.Movie

	stringQuery := masterQueries[GetList]

	param, err := json.Marshal(params)
	if err != nil {
		log.Println("marshal err: ", err)
		return nil, err
	}

	err = mr.redis.WithCache(ctx, fmt.Sprintf(GetListMoviesRedisKey, param), &Movie, func() (interface{}, error) {
		rows, err := mr.db.NamedQueryContext(ctx, stringQuery, params)
		if err != nil {
			log.Println("named query err: ", err)
			return nil, err
		}

		for rows.Next() {
			var dataMovie entity.Movie
			err = rows.StructScan(&dataMovie)
			if err != nil {
				return nil, err
			}

			Movie = append(Movie, &dataMovie)
		}

		return Movie, nil
	})

	if err != nil {
		log.Println("GetMovieList err: ", err)
		return nil, err
	}

	return Movie, nil
}

func (mr *MoviesRepository) GetMovieCount(ctx context.Context, param contract.GetListParam) (int64, error) {
	var count int64

	params, err := json.Marshal(param)
	if err != nil {
		log.Println("marshal err: ", err)
		return 0, err
	}

	err = mr.redis.WithCache(ctx, fmt.Sprintf(GetMoviesCountRedisKey, params), &count, func() (interface{}, error) {
		var countData int64
		err := mr.masterStmts[GetCountList].Get(&countData)
		return countData, err
	})

	if err != nil {
		log.Println("GetMovieCount err: ", err)
		return 0, err
	}

	return count, nil
}

func (mr *MoviesRepository) Get(ctx context.Context, id int) (entity.Movie, error) {
	var Movie entity.Movie
	err := mr.redis.WithCache(ctx, fmt.Sprintf(GetDetailMoviesRedisKey, id), &Movie, func() (interface{}, error) {
		var MovieData entity.Movie
		err := mr.masterStmts[GetByID].GetContext(ctx, &MovieData, id)
		return MovieData, err
	})

	if err != nil {
		log.Println(err)
		return Movie, err
	}

	return Movie, nil
}

func (mr *MoviesRepository) Update(ctx context.Context, data *entity.Movie) error {
	var rowsAffected int64

	namedStmt, err := mr.getNamedStatement(ctx, UpdateMovie)
	if err != nil {
		log.Println("get named statement err: ", err)
		return err
	}

	res, err := namedStmt.ExecContext(ctx, data)
	if err != nil {
		log.Println("exec err: ", err)
		return err
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		log.Println("Get rows affected err: ", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("ID not exist err: ", sql.ErrNoRows)
		return sql.ErrNoRows
	}

	redisErr := mr.redis.DelWithPattern(ctx, DeleteMovieRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return nil
}

func (mr *MoviesRepository) Create(ctx context.Context, data *entity.Movie) (contract.MovieResponseDB, error) {
	var res contract.MovieResponseDB

	namedStmt, err := mr.getNamedStatement(ctx, InsertMovie)
	if err != nil {
		log.Println("getNamedStatement err: ", err)
		return res, err
	}

	if err = namedStmt.GetContext(ctx, &res, data); err != nil {
		log.Println("get invoice err: ", err)
		return res, err
	}

	redisErr := mr.redis.DelWithPattern(ctx, DeleteMovieRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return res, nil
}

func (mr *MoviesRepository) Delete(ctx context.Context, id int64) error {
	stmt, err := mr.getStatement(ctx, Delete)
	if err != nil {
		log.Println("delete err: ", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		log.Println("delete err: ", err)
		return err
	}

	err = mr.redis.DelWithPattern(ctx, DeleteMovieRedisKey)
	if err != nil {
		log.Println("delete redis err: ", err)
	}

	return nil
}
