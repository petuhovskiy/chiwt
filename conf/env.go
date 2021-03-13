package conf

import "github.com/caarlos0/env/v6"

func ReadFromEnv() (*App, error) {
	var cfg App
	err := env.Parse(&cfg)
	return &cfg, err
}
