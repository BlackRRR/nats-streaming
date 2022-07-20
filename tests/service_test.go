package tests

import (
	nats_streaming "github.com/BlackRRR/nats-streaming/internal/app/services/nats-streaming"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"github.com/patrickmn/go-cache"
	"testing"
)

func Test_ServiceSave(t *testing.T) {
	sub := &stan.Msg{
		MsgProto: pb.MsgProto{
			Data: []byte(`{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`),
		},
	}
	service := nats_streaming.Service{
		Repositories: nil,
		Cache:        cache.New(0, 0),
	}

	err := service.SaveModelInPostgres(sub)
	if err != nil {
		t.Errorf("failed to save %s", err.Error())
	}

	order, ok := service.Cache.Get("b563feb7b2b84b6test")
	if order == nil || !ok {
		t.Error("wrong id")
	}

	t.Log(order, ok)
}

func Test_ServiceGet(t *testing.T) {
	id := "b563feb7b2b84b6test"
	service := nats_streaming.Service{
		Repositories: nil,
		Cache:        cache.New(0, 0),
	}

	order, err := service.GetOrder(id)
	if err != nil {
		t.Errorf("failed to get %s", err.Error())
	}

	t.Log(order)
}
