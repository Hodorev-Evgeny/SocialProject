package core_repository_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type PostgresConfig struct {
	Username string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" envDefault:"5432"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewPostgresConfig() (PostgresConfig, error) {
	var config PostgresConfig
	if envconfig.Process("POSTGRES", &config) != nil {
		return config, fmt.Errorf("error when try get environment variables")
	}

	return config, nil
}

func MustPostgresConfig() PostgresConfig {
	config, err := NewPostgresConfig()
	if err != nil {
		panic(err)
	}

	return config
}
