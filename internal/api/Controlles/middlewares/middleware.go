package middlewares

import (
	"Workmate/internal/service/authService"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type MiddlewaresInterface interface {
	CheckUsersToken(ctx *gin.Context)
}

type Middlewares struct {
	Auth authService.AuthServiceInterface
}

func NewMiddleware(auth authService.AuthServiceInterface) MiddlewaresInterface {
	return &Middlewares{Auth: auth}
}

func (m Middlewares) CheckUsersToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	token = strings.Split(token, " ")[1]

	if token == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	userId, err := m.Auth.ValidateToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	}
	ctx.Set("id", userId)
	ctx.Next()
}
