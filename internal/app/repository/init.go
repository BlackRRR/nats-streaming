package repository

import (
	"context"
	nats_streaming "github.com/BlackRRR/nats-streaming/internal/app/repository/nats-streaming"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repositories struct {
	*nats_streaming.Repository
}

func InitRepositories(ctx context.Context, pool *pgxpool.Pool) (*Repositories, error) {
	repository, err := nats_streaming.NewRepository(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get new postgres repository")
	}

	return &Repositories{
		repository,
	}, nil
}
