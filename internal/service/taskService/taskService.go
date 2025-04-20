package taskService

import (
	db "Workmate/internal/repositories/postgresql/sqlc"
	"Workmate/internal/repositories/redis"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
	GetTaskInfo(ctx context.Context, TaskID string) (db.Task, error)
}

type TaskService struct {
	querier  db.Querier
	Redis    redis.RedisRepositoryInterface
	WaitTime int
}

func NewTaskService(querier db.Querier, repositoryInterface redis.RedisRepositoryInterface, time int) TaskServiceInterface {
	return &TaskService{
		querier:  querier,
		Redis:    repositoryInterface,
		WaitTime: time,
	}
}
func (t *TaskService) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	taskId := uuid.New()
	arg.ID = taskId.String()
	arg.Status = "created"
	//Сначала заносим задачу в бд на долгосрочную перспективу
	task, err := t.querier.CreateTask(ctx, arg)
	if err != nil {
		return db.Task{}, fmt.Errorf("cannot create task %v", err)
	}

	//оставим задачу выполняться
	go t.DoSomethingWithTask(ctx, task)
	return task, nil
}

func (t *TaskService) GetTaskInfo(ctx context.Context, TaskID string) (db.Task, error) {
	//Либо получаем задачу из кеша - в первые 5 минут
	task, err := t.Redis.GetTaskData(ctx, TaskID)
	if err != nil {
		log.Printf("cannot find task in cache %v\n", err)
	} else {
		return task, nil
	}
	//Либо получаем задачу из бд - в любое другое время
	task, err = t.querier.GetTask(ctx, TaskID)
	if err != nil {
		return db.Task{}, fmt.Errorf("task not found %v", err)
	}
	return task, nil

}

// Будем генерировать статус и результат на основе рандома

var statusMap = map[bool]string{
	true:  "сompleted with an error",    //была ошибка(искусственно)
	false: "task successfully finished", // не было ошибки
}
var resultMap = map[bool]string{
	true:  "no result because of error",  // была ошибка
	false: "result of task : *somedata*", // не было ошибки
}

// искусственная нагрузка
func (t *TaskService) DoSomethingWithTask(ctx context.Context, task db.Task) {
	task.Status = "processing"
	//короткий кэш задачи на время исполнения операций
	err := t.Redis.CreateTaskWithData(ctx, 5*time.Minute, task)
	if err != nil {
		log.Printf("cannot cache task %v\n", err)
		return
	}
	// имитация действий над задачей
	//время выполнения задачи минут
	//randomTime := rand.Intn(2) + 1
	//time.Sleep(time.Duration(randomTime) * time.Minute)
	time.Sleep(time.Duration(t.WaitTime) * time.Second)
	//  имитация возникновения ошибки в процессе выполнения кода
	var WasErr bool
	r := rand.Intn(10)
	if r < 5 {
		WasErr = true
	} else {
		WasErr = false
	}
	fmt.Println(WasErr)
	statusParams := db.UpdateTaskStatusParams{
		ID:     task.ID,
		Status: statusMap[WasErr],
	}

	err = t.querier.UpdateTaskStatus(ctx, statusParams)
	if err != nil {
		log.Printf("cannot update task status %v", err)
		return
	}
	resultParams := db.UpdateTaskResultParams{
		ID:     task.ID,
		Result: resultMap[WasErr],
	}
	err = t.querier.UpdateTaskResult(ctx, resultParams)
	if err != nil {
		log.Printf("cannot updated task result %v", err)
		return
	}
	task, err = t.querier.GetTask(ctx, task.ID)
	if err != nil {
		log.Printf("cannot find task %v", err)
	}
	//кеш задачи на более долгое время чтоб быстрее выдавать пользователю информацию по результатам его задаче
	err = t.Redis.CreateTaskWithData(ctx, 2*time.Hour, task)
	if err != nil {
		log.Printf("cannot cache task %v\n", err)
		return
	}
}
