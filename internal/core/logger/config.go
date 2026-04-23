package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

/*
создания стурктуры для получения настроек логера
получения настроект происходит с помощью окружающей среды
*/

type Config struct {
	// Create setting for logger
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

/*
функции для создания логера
первая функция пытаеться создать конфиг и при неудаче возврощает соответстыущую ошибку
а уже в маст в случае ошибки кидает панику
*/

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("failed to process env config: %w", err)
	}

	return config, nil
}

func MustNewConfig() Config {
	config, err := NewConfig()

	if err != nil {
		err := fmt.Errorf("failed to create config: %w", err)
		panic(err)
	}

	return config
}
