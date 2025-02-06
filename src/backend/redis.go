package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rip/lib/config"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// InitializeRedis инициализирует клиент Redis
func InitializeRedis() (*redis.Client, error) {
	addr, password, db, err := config.FromEnvRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis configuration from environment: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	// Очищаем хранилище (текущую базу данных)
	err = client.FlushDB(ctx).Err() // или client.FlushAll(ctx) для очистки всех баз данных
	if err != nil {
		return nil, fmt.Errorf("failed to flush Redis database: %v", err)
	}

	log.Println("Redis client initialized successfully")
	return client, nil
}

type Session struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// SaveSession сохраняет сессию пользователя в Redis
func SaveSession(ctx context.Context, redisClient *redis.Client, userID uuid.UUID, role string, expiration time.Duration) error {
	session := Session{
		UserID: userID.String(),
		Role:   role,
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = redisClient.Set(ctx, userID.String(), sessionData, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// CheckSessionExists проверяет, существует ли сессия пользователя
func CheckSessionExists(ctx context.Context, redisClient *redis.Client, userID string) (bool, error) {
	exists, err := redisClient.Exists(ctx, userID).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// GetSession получает сессию пользователя из Redis
func GetSession(ctx context.Context, redisClient *redis.Client, userID string) (*Session, error) {
	sessionData, err := redisClient.Get(ctx, userID).Result()
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession удаляет сессию пользователя из Redis
func DeleteSession(ctx context.Context, redisClient *redis.Client, userID string) error {
	return redisClient.Del(ctx, userID).Err()
}
