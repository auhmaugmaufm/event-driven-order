package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	RDBHost string
	RDBPort string

	JWTSecret     string
	JWTExpireHour int
}

var cfg *Config

func Load() {
	if os.Getenv("APP_NEW") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("no .env file found")
		}
	}

	expireHour, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOUR", "24"))
	cfg = &Config{
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),

		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		DBSSLMode: os.Getenv("DB_SSLMODE"),

		RDBHost: os.Getenv("RDB_HOST"),
		RDBPort: os.Getenv("RDB_PORT"),

		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpireHour: expireHour,
	}
}

func Get() *Config {
	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
