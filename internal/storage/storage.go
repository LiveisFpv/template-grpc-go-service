package storage

import (
	"context"
	"fmt"
	"template-grpc-go-service/internal/domain/models"
	postgresql "template-grpc-go-service/internal/storage/postgreSQL"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type repo struct {
	*postgresql.Queries
	log  *logrus.Logger
	pool *pgxpool.Pool
}

// Repository constructor
func NewRepository(
	pgxpool *pgxpool.Pool,
	log *logrus.Logger,
) Repository {
	return &repo{
		Queries: postgresql.New(pgxpool),
		log:     log,
		pool:    pgxpool,
	}
}

// Func for work with DB
type Repository interface {
	//Func needed for DB
	GetNamebyID(ctx context.Context, name_id int) (name *models.Name, err error)
	GetAllName(ctx context.Context, pagination *models.Pagination, filter []*models.Filter, orderby []*models.Sort) ([]*models.Name, *models.Pagination, error)
	CreateName(ctx context.Context, name_title, name_capital, name_area string) (name *models.Name, err error)
	UpdateNamebyID(ctx context.Context, name *models.Name) (err error)
	DeleteNamebyID(ctx context.Context, name_id int) (name *models.Name, err error)
	Stop()
}

// Создание нового подключения к БД
func NewStorage(ctx context.Context, dsn string, log *logrus.Logger) (Repository, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return NewRepository(pool, log), nil
}

func (r *repo) Stop() {
	r.Queries.Stop()
}
