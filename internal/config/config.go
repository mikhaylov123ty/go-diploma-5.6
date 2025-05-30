package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Address              string
	DBURI                string
	AccuralSystemAddress string
	Secret               string
	LogLevel             string
}

func Init() (*Config, error) {
	res := &Config{}

	res.parseFlags()

	if err := res.parseEnvironment(); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	//TODO change it to something more interesting
	host, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("error getting hostname: %w", err)
	}
	res.Secret = host

	return res, nil
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.Address, "a", ":8080", "Address to listen on. Default: \":8080\"")

	flag.StringVar(&c.DBURI, "d", "", "Database URI: 'postgresql://postgres:postgres@hostname/postgres?sslmode=disable'")

	flag.StringVar(&c.AccuralSystemAddress, "r", "", "Address of accural system: http://hostname:port")

	flag.StringVar(&c.LogLevel, "l", "info", "Log level. Default: info")

	flag.Parse()
}

func (c *Config) parseEnvironment() error {
	if os.Getenv("RUN_ADDRESS") != "" {
		c.Address = os.Getenv("RUN_ADDRESS")
	}

	if os.Getenv("DATABASE_URI") != "" {
		c.DBURI = os.Getenv("DATABASE_URI")
	}

	if os.Getenv("ACCRUAL_SYSTEM_ADDRESS") != "" {
		c.AccuralSystemAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	}

	if os.Getenv("LOG_LEVEL") != "" {
		c.LogLevel = os.Getenv("LOG_LEVEL")
	}

	return nil
}
