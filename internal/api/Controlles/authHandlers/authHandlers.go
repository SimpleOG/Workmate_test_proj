package authHandlers

import (
	db "Workmate/internal/repositories/postgresql/sqlc"
	"Workmate/internal/service/authService"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandlersInterface interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type AuthHandlers struct {
	authService.AuthServiceInterface
}

func NewAuthHandlers(serviceInterface authService.AuthServiceInterface) AuthHandlersInterface {
	return &AuthHandlers{
		serviceInterface,
	}
}
func (a *AuthHandlers) Register(ctx *gin.Context) {
	var userParams db.CreateUserParams
	if err := ctx.ShouldBindJSON(&userParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error : ": err})
		return
	}
	user, err := a.AuthServiceInterface.Register(userParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"User data": user})
}
func (a *AuthHandlers) Login(ctx *gin.Context) {
	var searchParams db.GetUserForLoginParams
	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error : ": err.Error()})
		return
	}
	token, err := a.AuthServiceInterface.Login(searchParams)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error : ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
