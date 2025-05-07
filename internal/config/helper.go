package config

import (
	"fmt"
	"time"
)

func (c *Config) GetReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeoutSec) * time.Second
}

func (c *Config) GetWriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeoutSec) * time.Second
}

func (c *Config) GetIdleTimeout() time.Duration {
	return time.Duration(c.IdleTimeoutSec) * time.Second
}

func (c *Config) GetConnMaxLifetime() time.Duration {
	return time.Duration(c.DBConnMaxLifetimeSec) * time.Second
}

// Data Source Name for postgres
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}
