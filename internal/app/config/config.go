package config

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ServicePort string          `yaml:"service_port"`
	PGConfig    *pgxpool.Config `yaml:"pg_config"`
}

func InitConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath("config")
	v.SetConfigName("config")

	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	var config Config

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?pool_max_conns=%s",
		v.Get("config.db_conn_config.user"),
		v.Get("config.db_conn_config.host"),
		v.Get("config.db_conn_config.password"),
		v.Get("config.db_conn_config.db_name"),
		v.Get("config.db_conn_config.max_pool_conns"))

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "`Init config` failed to parse config")
	}

	servicePort := v.GetString("config.service_port")

	config.PGConfig = pgxConfig
	config.ServicePort = servicePort

	return &config, nil
}
