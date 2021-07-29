package database

import "fmt"

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Instance string
}

// DSN generate data source name
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Instance)
}
