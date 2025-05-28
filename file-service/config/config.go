package config

// ::::::::::::::::::::::::::::::::::::::::::::::
// :::  Create own.env file for file-serivce  :::
// ::::::::::::::::::::::::::::::::::::::::::::::

import (
	"log"
	"os"
)

// Default values
const (
    POSTGRES_HOST	= "localhost"
    POSTGRES_PORT 	= "5432"
    POSTGRES_USER	= "postgres"
    POSTGRES_PASSWORD	= "postgres"
    POSTGRES_NAME	= "postgres"

    AUTH_TO_FILE_TOKEN	= "none"

    SERVER_PORT		= ":8082"
)

type Config struct {
    ServerPort string
    AuthToFileToken string
    
    PostgresConfig *postgresConfig
}

func NewConfig() *Config {
    if os.Getenv("AUTH_TO_FILE_TOKEN") == "" {
	os.Setenv("AUTH_TO_FILE_TOKEN", AUTH_TO_FILE_TOKEN)
    }

    config := &Config {
	ServerPort: SERVER_PORT,
	AuthToFileToken: os.Getenv("AUTH_TO_FILE_TOKEN"),
    }

    return config
}

func (c *Config) WithPostgres() *Config {
    if os.Getenv("POSTGRES_HOST") == "" {
	 os.Setenv("POSTGRES_HOST", POSTGRES_HOST)
    }
    if os.Getenv("POSTGRES_PORT") == "" {
	os.Setenv("POSTGRES_PORT", POSTGRES_PORT)
    }
    if os.Getenv("POSTGRES_USER") == "" {
	os.Setenv("POSTGRES_USER", POSTGRES_USER)
    }
    if os.Getenv("POSTGRES_PASSWORD") == "" {
	os.Setenv("POSTGRES_PASSWORD", POSTGRES_PASSWORD)
    }
    if os.Getenv("POSTGRES_NAME") == "" {
	os.Setenv("POSTGRES_NAME", POSTGRES_NAME)
    }

    c.PostgresConfig = &postgresConfig {
	Host: os.Getenv("POSTGRES_HOST"),
	Port: os.Getenv("POSTGRES_PORT"),
	User: os.Getenv("POSTGRES_USER"),
	Password: os.Getenv("POSTGRES_PASSWORD"),
	Name: os.Getenv("POSTGRES_NAME"),
    }

    log.Println("Config with Postgres initialized")
    return c
}
