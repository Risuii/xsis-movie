package movie

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	frsAtomic "github.com/Risuii/frs-lib/atomic"
	atomicSqlx "github.com/Risuii/frs-lib/atomic/sqlx"
	frsRedis "github.com/Risuii/frs-lib/redis"
	sqlxUtils "github.com/Risuii/frs-lib/sqlx"
)

const (
	AllFields = `id, title, description, rating, image, created_at, updated_at`

	GetByID = iota + 100
	GetByMovieID
	GetList
	GetCountList
	GetLatestMovieID
	Delete

	InsertMovie = iota + 200
	UpdateMovie

	// Redis Key

	GetListMoviesRedisKey   = "movie:movies:getlist:%s"
	GetDetailMoviesRedisKey = "movie:movies:getdetail:%d"
	GetMoviesCountRedisKey  = "movie:movies:getcount:%s"
	DeleteMovieRedisKey     = "movie:movies:*"
)

var (
	masterQueries = []string{
		GetByID:          fmt.Sprintf("SELECT %s FROM movies WHERE id = $1 AND deleted_at IS NULL", AllFields),
		GetByMovieID:     fmt.Sprintf("SELECT %s FROM movies WHERE id = $1 And deleted_at IS NULL", AllFields),
		GetList:          fmt.Sprintf(`SELECT %s FROM movies WHERE deleted_at IS NULL`, AllFields),
		GetCountList:     `SELECT COUNT(*) FROM movies WHERE deleted_at IS NULL`,
		GetLatestMovieID: `SELECT MAX(id) FROM movies`,
		Delete:           `UPDATE movies set deleted_at=now() WHERE id = $1`,
	}

	masterNamedQueries = []string{
		InsertMovie: `INSERT INTO movies (title, description, rating, image, created_at) VALUES (:title, :description, :rating, :image, now()) RETURNING id, title, description, rating, image, created_at, updated_at`,
		UpdateMovie: `UPDATE movies SET (title, description, rating, image, updated_at) = (:title, :description, :rating, :image, now()) WHERE id = :id`,
	}
)

type MoviesRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
	redis             frsRedis.Redis
}

func InitMoviesRepository(ctx context.Context, db *sqlx.DB, redis frsRedis.Redis) (*MoviesRepository, error) {
	stmpts, err := sqlxUtils.PrepareQueries(db, masterQueries)
	if err != nil {
		log.Println("PrepareQueries err:", err)
		return nil, err
	}

	namedStmpts, err := sqlxUtils.PrepareNamedQueries(db, masterNamedQueries)
	if err != nil {
		log.Println("PrepareNamedQueries err:", err)
		return nil, err
	}

	return &MoviesRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
		redis:             redis,
	}, nil
}

func (r *MoviesRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*frsAtomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = frsAtomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = r.masterStmts[queryId]
	}
	return statement, err
}

func (r *MoviesRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*frsAtomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = frsAtomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = r.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
