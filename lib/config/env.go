package config

import (
	"fmt"
	"os"
	"strconv"
)

// FromEnv собирает DSN строку из переменных окружения
func FromEnvDB() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

// FromEnvMinIO собирает настройки подключения для MinIO из переменных окружения
func FromEnvMinIO() (string, string, string, bool, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSLStr := os.Getenv("MINIO_USE_SSL")

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return "", "", "", false, fmt.Errorf("not all required environment variables are set")
	}

	useSSL := false
	if useSSLStr != "" {
		var err error
		useSSL, err = strconv.ParseBool(useSSLStr)
		if err != nil {
			return "", "", "", false, fmt.Errorf("could not convert MINIO_USE_SSL to bool: %v", err)
		}
	}

	return endpoint, accessKey, secretKey, useSSL, nil
}

// FromEnvRedis собирает настройки подключения к Redis из переменных окружения
func FromEnvRedis() (string, string, int, error) {
	addr := os.Getenv("REDIS_ADDR")     // Адрес Redis, например, "localhost:6379"
	password := os.Getenv("REDIS_PASS") // Пароль (может быть пустым)
	dbStr := os.Getenv("REDIS_DB")      // Номер базы данных в виде строки (например, "0")

	if addr == "" {
		return "", "", 0, fmt.Errorf("REDIS_ADDR is not set")
	}

	db := 0
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			return "", "", 0, fmt.Errorf("could not convert REDIS_DB to int: %v", err)
		}
	}

	return addr, password, db, nil
}

func LoadSecret() (string, error) {
	secret := os.Getenv("APP_SECRET")
	if secret == "" {
		return "", fmt.Errorf("missing APP_SECRET environment variable")
	}
	return secret, nil
}
