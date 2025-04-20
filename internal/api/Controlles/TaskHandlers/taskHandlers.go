package TaskHandlers

import (
	db "Workmate/internal/repositories/postgresql/sqlc"
	"Workmate/internal/service/taskService"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskHandlersInterface interface {
	CreateTask(ctx *gin.Context)
	GetTask(ctx *gin.Context)
}

type TaskHandlers struct {
	task taskService.TaskServiceInterface
}

func NewTaskHandlers(serviceInterface taskService.TaskServiceInterface) TaskHandlersInterface {
	return &TaskHandlers{task: serviceInterface}
}
func (t *TaskHandlers) CreateTask(ctx *gin.Context) {
	var arg db.CreateTaskParams
	task, err := t.task.CreateTask(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error : ": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"task": task})
}
func (t *TaskHandlers) GetTask(ctx *gin.Context) {
	taskId := ctx.Query("task_id")
	info, err := t.task.GetTaskInfo(ctx, taskId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error : ": err.Error()})
	}
	ctx.JSON(http.StatusOK, info)
}
