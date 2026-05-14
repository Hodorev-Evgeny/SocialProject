package mailresend

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config loads RESEND_* variables. API key is optional: if empty, OTP is not sent via Resend.
type Config struct {
	APIKey string `envconfig:"API_KEY"`
	From   string `envconfig:"FROM"`
}

func MustConfig() Config {
	var c Config
	if err := envconfig.Process("RESEND", &c); err != nil {
		panic(fmt.Errorf("resend config: %w", err))
	}
	return c
}
