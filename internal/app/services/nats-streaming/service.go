package nats_streaming

import (
	"encoding/json"
	"github.com/BlackRRR/nats-streaming/internal/app/repository"
	"github.com/BlackRRR/nats-streaming/internal/model"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"time"
)

type Service struct {
	*repository.Repositories
	*cache.Cache
}

func InitService(repo *repository.Repositories, cache *cache.Cache) *Service {
	return &Service{repo, cache}
}

func (s *Service) SaveModelInPostgres(msg *stan.Msg) error {
	var data *model.Model
	var pgModel model.PostgresModel

	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Service: failed to unmarshal order")
	}

	//Check for correct data comes
	if data == nil {
		return errors.Wrap(err, "Service: wrong data from message")
	}

	if data.OrderUid == "" {
		return errors.Wrap(err, "Service: wrong data from message")
	}

	pgModel.Id = data.OrderUid
	pgModel.Order, err = postgresqlModel(data)
	if err != nil {
		return errors.Wrap(err, "Service: failed marshal postgresModel")
	}

	//save in postgres
	err = s.PostgresRepository.CreateOrder(pgModel.Id, pgModel.Order)
	if err != nil {
		return errors.Wrap(err, "Service: failed to create order")
	}

	return nil
}

func (s *Service) SaveModelInCache(msg *stan.Msg) error {
	var data model.Model

	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Service: failed to unmarshal order")
	}

	//Check for correct data comes
	if data.OrderUid == "" {
		return errors.Wrap(err, "Service: wrong data from message")
	}

	//save in memory
	err = s.Cache.Add(data.OrderUid, data, time.Hour*100)
	if err != nil {
		return errors.Wrap(err, "Service: failed to add message in cache")
	}

	return nil
}

func (s *Service) DataRecovery() error {
	if s.Cache.ItemCount() == 0 {
		orders, err := s.PostgresRepository.GetOrders()
		if err != nil {
			return errors.Wrap(err, "Service: failed to get orders")
		}

		if orders == nil {
			return nil
		}

		for i := range orders {
			err := s.Cache.Add(orders[i].Id, orders[i].Order, time.Hour*100)
			if err != nil {
				return errors.Wrap(err, "Service: failed to add order in cache from postgres")
			}
		}
	}

	return nil
}

func (s *Service) GetOrder(id string) (json.RawMessage, error) {
	orderFromCache, indicator := s.Cache.Get(id)

	if !indicator {
		return nil, nil
	}

	order, err := json.MarshalIndent(orderFromCache, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "Service: failed convert order from cache into json")
	}

	return order, nil
}

func postgresqlModel(data *model.Model) (json.RawMessage, error) {
	order, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return order, nil
}
