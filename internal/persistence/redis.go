package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"poker-server/internal/lobby"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) (*RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisStore{client: rdb}, nil
}

func (r *RedisStore) SaveLobby(l *lobby.Lobby) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Locking NOT required here if caller handles concurrency (Manager does).
	// We just serialize what is passed.

	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("failed to marshal lobby: %w", err)
	}

	key := fmt.Sprintf("poker:lobby:%s", l.Code)
	// No expiration for now, or long expiration (e.g. 24h)
	if err := r.client.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	return nil
}

func (r *RedisStore) GetLobby(code string) (*lobby.Lobby, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("poker:lobby:%s", code)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("redis get failed: %w", err)
	}

	var l lobby.Lobby
	if err := json.Unmarshal([]byte(data), &l); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lobby: %w", err)
	}
	// Restore internal references if needed (e.g. references to config or something not serialized)
	// For MVP, simple structs usually fine.

	// Key: The deserialized lobby needs to be valid.
	// If GameTable is present, its internal state must be consistent.
	// Since we serialized the whole struct tree, it should be fine.

	return &l, nil
}
