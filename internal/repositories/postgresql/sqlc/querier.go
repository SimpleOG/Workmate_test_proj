// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteTask(ctx context.Context, id string) error
	GetTask(ctx context.Context, id string) (Task, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserForLogin(ctx context.Context, arg GetUserForLoginParams) (User, error)
	UpdateTaskResult(ctx context.Context, arg UpdateTaskResultParams) error
	UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error
}

var _ Querier = (*Queries)(nil)
