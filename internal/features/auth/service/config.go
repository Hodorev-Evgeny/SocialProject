package features_auth_service

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type JWTConfig struct {
	Secret     string        `envconfig:"SECRET" required:"true"`
	AccessTTL  time.Duration `envconfig:"ACCESS_TTL" envDefault:"15m"`
	RefreshTTL time.Duration `envconfig:"REFRESH_TTL" envDefault:"720h"`
}

func MustJWTConfig() JWTConfig {
	var c JWTConfig
	if err := envconfig.Process("JWT", &c); err != nil {
		panic(fmt.Errorf("jwt config: %w", err))
	}
	if len(c.Secret) < 32 {
		panic("JWT_SECRET must be at least 32 characters")
	}
	return c
}
