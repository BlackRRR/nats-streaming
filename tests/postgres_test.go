package tests

import (
	"context"
	"github.com/BlackRRR/nats-streaming/internal/app/config"
	nats_streaming "github.com/BlackRRR/nats-streaming/internal/app/repository/nats-streaming"
	"github.com/jackc/pgx/v4/pgxpool"
	"testing"
)

func TestPostgres(t *testing.T) {
	ctx := context.Background()

	initConfig, err := config.InitConfig()
	if err != nil {
		t.Errorf("failed to init config %s", err.Error())
	}

	dbConn, err := pgxpool.ConnectConfig(ctx, initConfig.PGConfig)
	if err != nil {
		t.Errorf("failed to init postgres %s", err.Error())
	}

	postgresRepository := &nats_streaming.PostgresRepository{Ctx: ctx, ConnPool: dbConn}

	err = postgresRepository.CreateOrder("b563feb7b2b84b6test", []byte(`"order_uid": "b563feb7b2b84b6test",
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
	}`))
	if err != nil {
		t.Errorf("failed to add order %s", err.Error())
	}

	order, err := postgresRepository.GetOrder("b563feb7b2b84b6test")
	if err != nil {
		t.Errorf("failed to get order %s", err.Error())
	}

	t.Log(order)
}
