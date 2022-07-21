package nats_streaming

import (
	"context"
	"encoding/json"
	"github.com/BlackRRR/nats-streaming/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	Ctx context.Context

	ConnPool *pgxpool.Pool
}

func NewRepository(ctx context.Context, connPool *pgxpool.Pool) (*Repository, error) {
	pgRep := &Repository{Ctx: ctx, ConnPool: connPool}

	_, err := pgRep.ConnPool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS orders(
id text UNIQUE,
order_model json);`)
	if err != nil {
		return nil, errors.Wrap(err, "Postgres: failed to create table")
	}
	return pgRep, nil
}

func (r *Repository) CreateOrder(id string, order json.RawMessage) error {
	_, err := r.ConnPool.Exec(r.Ctx, `
INSERT INTO orders 
(id, order_model) 
VALUES ($1,$2)`, id, order)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"orders_id_key\" (SQLSTATE 23505)" {
			return nil
		}
		return errors.Wrap(err, "Postgres: failed to insert into postgres")
	}

	return nil
}

func (r *Repository) GetOrder(id string) (json.RawMessage, error) {
	var data json.RawMessage
	err := r.ConnPool.QueryRow(r.Ctx, `SELECT order_model FROM orders WHERE id = $1`, id).Scan(&data)
	if err != nil {
		return nil, errors.Wrap(err, "Postgres: failed to get order from database")
	}

	return data, nil
}

func (r *Repository) GetOrders() ([]*model.RepositoryModel, error) {
	rows, err := r.ConnPool.Query(r.Ctx, `SELECT * FROM orders`)
	if err != nil {
		return nil, errors.Wrap(err, "Postgres: failed to get order from database")
	}

	orders, err := ReadRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get orders from rows")
	}

	return orders, nil
}

func ReadRows(rows pgx.Rows) ([]*model.RepositoryModel, error) {
	order := &model.RepositoryModel{}
	var orders []*model.RepositoryModel

	for rows.Next() {
		err := rows.Scan(&order.Id, &order.Order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
