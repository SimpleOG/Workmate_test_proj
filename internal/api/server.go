package api

import (
	"Workmate/internal/api/Controlles/TaskHandlers"
	"Workmate/internal/api/Controlles/authHandlers"
	"Workmate/internal/api/Controlles/middlewares"

	"Workmate/internal/service"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	authHandlers.AuthHandlersInterface
	TaskHandlers.TaskHandlersInterface
	middlewares.MiddlewaresInterface
}

func NewControllers(services service.Services) Controllers {
	return Controllers{
		AuthHandlersInterface: authHandlers.NewAuthHandlers(services.AuthServiceInterface),
		TaskHandlersInterface: TaskHandlers.NewTaskHandlers(services.TaskServiceInterface),
		MiddlewaresInterface:  middlewares.NewMiddleware(services.AuthServiceInterface),
	}
}

type Server struct {
	Router      *gin.Engine
	Controllers Controllers
}

func NewServer(router *gin.Engine, services service.Services) Server {
	return Server{
		Router:      router,
		Controllers: NewControllers(services),
	}
}
func (s *Server) Run(address string) error {
	s.InitRoutes()
	return s.Router.Run(address)
}
func (s *Server) InitRoutes() {
	s.InitAuthRoutes()
	s.InitTaskRoutes()
}
func (s *Server) InitAuthRoutes() {
	auth := s.Router.Group("/auth")
	{
		auth.POST("/sign_in", s.Controllers.AuthHandlersInterface.Register)
		auth.POST("/login", s.Controllers.AuthHandlersInterface.Login)
	}
}
func (s *Server) InitTaskRoutes() {
	task := s.Router.Group("/task", s.Controllers.MiddlewaresInterface.CheckUsersToken)
	{
		task.POST("/create_task", s.Controllers.TaskHandlersInterface.CreateTask)
		task.GET("/get_task", s.Controllers.TaskHandlersInterface.GetTask)
	}
}
