package redis

import (
	db "Workmate/internal/repositories/postgresql/sqlc"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepositoryInterface interface {
	CreateTaskWithData(ctx context.Context, expTime time.Duration, task db.Task) error
	GetTaskData(ctx context.Context, taskID string) (db.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
}
type RedisRepository struct {
	RedisClient *redis.Client
	ExpTime     time.Duration
}

func NewRedisClient(addr string) (RedisRepositoryInterface, error) {

	client := redis.NewClient(&redis.Options{Addr: addr, Password: "", DB: 0})
	repo := &RedisRepository{RedisClient: client}
	_, err := repo.RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return repo, err
}

func (r *RedisRepository) CreateTaskWithData(ctx context.Context, expTime time.Duration, task db.Task) error {
	taskKey := fmt.Sprintf("task: %v", task.ID)
	value, err := json.Marshal(task)
	if err != nil {
		fmt.Errorf("cannot marshal task data")
	}
	return r.RedisClient.Set(ctx, taskKey, value, expTime).Err()
}
func (r *RedisRepository) GetTaskData(ctx context.Context, taskID string) (db.Task, error) {
	taskKey := fmt.Sprintf("task: %v", taskID)

	res, err := r.RedisClient.Get(ctx, taskKey).Result()
	if err != nil {
		return db.Task{}, fmt.Errorf("cannot find task %v", err)
	}
	var task db.Task
	if err = json.Unmarshal([]byte(res), &task); err != nil {
		return db.Task{}, fmt.Errorf("cannot parse task data %v", err)
	}

	return task, nil
}

func (r *RedisRepository) DeleteTask(ctx context.Context, taskID string) error {
	key := fmt.Sprintf("task: %v", taskID)
	return r.RedisClient.Del(ctx, key).Err()
}
