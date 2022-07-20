package main

import (
	"context"
	"github.com/BlackRRR/nats-streaming/internal/app/config"
	"github.com/BlackRRR/nats-streaming/internal/app/repository"
	"github.com/BlackRRR/nats-streaming/internal/app/server"
	nats_streaming "github.com/BlackRRR/nats-streaming/internal/app/services/nats-streaming"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.TODO()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("failed to init config %s", err.Error())
	}

	log.Println("init config success")

	dbConn, err := pgxpool.ConnectConfig(ctx, cfg.PGConfig)
	if err != nil {
		log.Fatalf("failed to init postgres %s", err.Error())
	}

	log.Println("connect database success")

	repositories, err := repository.InitRepositories(ctx, dbConn)
	if err != nil {
		log.Fatalf("failed to init repository %s", err.Error())
	}

	log.Println("Init repository success")

	newCache := cache.New(0, 0)
	service := nats_streaming.InitService(repositories, newCache)

	err = service.DataRecovery()
	if err != nil {
		log.Fatalf("failed to recovery cache from postgres %s", err.Error())
	}

	log.Println("Init service success")

	sc, err := stan.Connect("test-cluster", "myID",
		stan.Pings(1, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.", err)
	}

	log.Println("Connected NATS Streaming subsystem")

	handler := server.NewHandler(service)

	log.Println("Init handler success")

	go func() {
		log.Println("http server started on port:", cfg.ServicePort)
		serviceErr := http.ListenAndServe(":"+cfg.ServicePort, handler.InitRouters())
		if serviceErr != nil {
			log.Fatalf("http handler was stoped by err: %s", serviceErr.Error())
		}
	}()

	go func(stan.Conn) {
		_, err = sc.Subscribe("foo", func(msg *stan.Msg) {
			//fmt.Printf("new data %s", string(msg.Data))
			if err != nil {
				log.Printf("failed to subscribe %s", err.Error())
			}

			err = service.SaveModelInPostgres(msg)
			if err != nil {
				log.Printf("service failed `Save model in postgres` %s", err.Error())
			}

			err = service.SaveModelInCache(msg)
			if err != nil {
				log.Printf("service failed `Save model in cache` %s", err.Error())
			}
		})
	}(sc)

	sig := <-subscribeToSystemSignals()

	log.Printf("shutdown all process on '%s' system signal\n", sig.String())
}

func subscribeToSystemSignals() chan os.Signal {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	return ch
}
