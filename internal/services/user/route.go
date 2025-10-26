package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	r := rg.Group("/users")
	r.POST("/register", handler.Register)
	r.POST("/auth/login", handler.Login)
}
