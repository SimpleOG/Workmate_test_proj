package service

import (
	db "Workmate/internal/repositories/postgresql/sqlc"
	"Workmate/internal/repositories/redis"
	"Workmate/internal/service/authService"
	"Workmate/internal/service/taskService"
	"Workmate/util/config"
)

type Services struct {
	authService.AuthServiceInterface
	taskService.TaskServiceInterface
}

func NewServices(querier db.Querier, repositoryInterface redis.RedisRepositoryInterface, config config.Config) Services {
	return Services{
		AuthServiceInterface: authService.NewAuthService(querier, config.SecretKey),
		TaskServiceInterface: taskService.NewTaskService(querier, repositoryInterface, config.WaitTime),
	}
}
