package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sportivaid/go-template/src/account"
	"github.com/sportivaid/go-template/src/common/auth"
)

func NewUserHandler(am *auth.Middleware, au account.Usecase) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/login", Login(au))
		v1.GET("/users", GetUsers(au))
		v1.GET("/user", am.AuthUserToken(GetUser(au)))
	}

	return router
}
