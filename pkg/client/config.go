package client

import (
	"fmt"
	"net"
)

type Config struct {
	User string
	Password string
	Host string
	Port uint16
	Database string
	hostIP net.IP
}

func NewConfig(host string, port uint16, user, pwd, db string) *Config {
	return &Config{
		user,
		pwd,
		host,
		port,
		db,
		nil,
	}
}

func (c *Config) Validate() error {
	if c.User == "" {
		c.User = "postgres"
	}

	if c.Host == "" {
		return fmt.Errorf(`"host" required`)
	}

	if c.Database == "" {
		return fmt.Errorf(`"Database" required`)
	}

	if c.Port == 0 {
		c.Port = 5432
	}

	c.hostIP = net.ParseIP(c.Host)
	if c.hostIP == nil {
		IPs, err := net.LookupIP(c.Host)
		if err != nil {
			return err
		}
		c.hostIP = IPs[0]
	}

	return nil
}