package config

import (
	"fmt"
)

type postgresConfig struct {
    Host	string
    Port	string
    User	string
    Password	string
    Name	string
}

func (c *postgresConfig) GetConnString() string {
    connString := fmt.Sprintf(
	"postgres://%s:%s@%s:%s/%d?sslmode=disable",
	c.User,
	c.Password,
	c.Host,
	c.Port,
	c.Name,
    )

    return connString
}
