package config

import (
	"os"
	"time"
)

type Config struct {
	Address         string
	StoreType       string
	SqliteDSN       string
	Timeout         time.Duration
}

func Load() (*Config, error) {
	return &Config{
		Address:         getEnv("MODBUS_SERVER_ADDRESS", ":502"),
		StoreType:       getEnv("MODBUS_SERVER_STORE_TYPE", "inmemory"),
		SqliteDSN:       getEnv("MODBUS_SERVER_SQLITE_DSN", "modbus.db"),
		Timeout:         func() time.Duration {
			value, exists := os.LookupEnv("MODBUS_SERVER_TIMEOUT")
			if exists {
				duration, err := time.ParseDuration(value)
				if err == nil {
					return duration
				}
			}
			return 5 * time.Second
		}(),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Implement getEnvAsDuration similarly...

