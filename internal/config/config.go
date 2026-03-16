package config

import (
	"cmp"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDB       string
	ServerAddress string
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQLUser, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDB)
}

func Load() (Config, error) {
	_ = godotenv.Load() // optional: env vars may already be set (e.g. Docker env_file)

	return Config{
		MySQLHost:     os.Getenv("MYSQL_HOST"),
		MySQLPort:     os.Getenv("MYSQL_PORT"),
		MySQLUser:     os.Getenv("MYSQL_USER"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
		MySQLDB:       os.Getenv("MYSQL_DB"),
		ServerAddress: cmp.Or(os.Getenv("SERVER_ADDRESS"), ":8080"),
	}, nil
}
