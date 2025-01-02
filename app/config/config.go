package config

import "encore.dev/config"

type Config struct {
	TemporalHostPort config.String
	DBHost           config.String
	DBPort           config.String
	DBUser           config.String
	DBPassword       config.String
	DBName           config.String
}
