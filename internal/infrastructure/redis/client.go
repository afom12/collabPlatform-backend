package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/collab-platform/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	pubsub *redis.PubSub
	ctx    context.Context
}

func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis successfully")
	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func (r *RedisClient) PublishOperation(docID string, message domain.BroadcastMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	channel := fmt.Sprintf("document:%s", docID)
	return r.client.Publish(r.ctx, channel, data).Err()
}

func (r *RedisClient) SubscribeToDocument(docID string) (*redis.PubSub, error) {
	channel := fmt.Sprintf("document:%s", docID)
	pubsub := r.client.Subscribe(r.ctx, channel)
	return pubsub, nil
}

func (r *RedisClient) SetUserSession(userID, docID string, data interface{}) error {
	key := fmt.Sprintf("session:%s:%s", userID, docID)
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, value, 0).Err()
}

func (r *RedisClient) GetUserSession(userID, docID string) (string, error) {
	key := fmt.Sprintf("session:%s:%s", userID, docID)
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Close() error {
	if r.pubsub != nil {
		r.pubsub.Close()
	}
	return r.client.Close()
}

