package nats_streaming

import (
	"encoding/json"
	"github.com/BlackRRR/nats-streaming/internal/app/repository"
	"github.com/BlackRRR/nats-streaming/internal/model"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type Service struct {
	*repository.Repositories
	*cache.Cache
}

func InitService(repo *repository.Repositories, cache *cache.Cache) *Service {
	return &Service{repo, cache}
}

func (s *Service) SaveModelInRepository(msg *stan.Msg) error {
	var data *model.Order
	var pgModel model.RepositoryModel

	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Service: failed to unmarshal order")
	}

	//Check for correct data comes
	reg := regexp.MustCompile("^[0-9a-z]+$")
	CheckUUID := reg.MatchString(data.OrderUid)

	if !CheckUUID {
		return errors.New("Service: wrong UUID")
	}

	if data == nil {
		return errors.Wrap(err, "Service: wrong data from message")
	}

	if data.OrderUid == "" {
		return errors.Wrap(err, "Service: wrong data from message")
	}

	pgModel.Id = data.OrderUid
	pgModel.Order, err = JsonRawModel(data)
	if err != nil {
		return errors.Wrap(err, "Service: failed marshal Order")
	}

	//save in repository
	err = s.Repository.CreateOrder(pgModel.Id, pgModel.Order)
	if err != nil {
		return errors.Wrap(err, "Service: failed to create order")
	}

	return nil
}

func (s *Service) SaveModelInCache(msg *stan.Msg) error {
	var data *model.Order

	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return errors.Wrap(err, "Service: failed to unmarshal order")
	}

	//Check for correct data comes
	reg := regexp.MustCompile("^[0-9a-z]+$")
	CheckUUID := reg.MatchString(data.OrderUid)

	if !CheckUUID {
		return errors.New("Service: wrong UUID")
	}

	if data == nil {
		return errors.Wrap(err, "Service: wrong data from message")
	}

	if data.OrderUid == "" {
		return errors.Wrap(err, "Service: wrong data from message")
	}

	//save in memory
	err = s.Cache.Add(data.OrderUid, data, cache.NoExpiration)
	if err != nil {
		exist := strings.Contains(err.Error(), "already exists")
		if exist {
			return nil
		}
		return errors.Wrap(err, "Service: failed to add message in cache")
	}

	return nil
}

func (s *Service) DataRecovery() error {
	//recovery data if cache null
	if s.Cache.ItemCount() == 0 {
		orders, err := s.Repository.GetOrders()
		if err != nil {
			return errors.Wrap(err, "Service: failed to get orders")
		}

		if orders == nil {
			return nil
		}

		for _, order := range orders {
			err := s.Cache.Add(order.Id, order.Order, cache.NoExpiration)
			if err != nil {
				exist := strings.Contains(err.Error(), "already exists")
				if exist {
					return nil
				}
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

func JsonRawModel(data *model.Order) (json.RawMessage, error) {
	order, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return order, nil
}
