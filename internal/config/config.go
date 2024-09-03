package config

import "fmt"

type App struct {
	Database Database
}

type Database struct {
	Host           string
	Port           int
	Name           string
	User           string
	Password       string
	MigrationsPath string
}

func (c Database) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable default_query_exec_mode=cache_describe",
		c.Host, c.Port, c.User, c.Name, c.Password)
}
