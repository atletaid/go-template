package auth

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sportivaid/go-template/src/common/apperror"
	"github.com/sportivaid/go-template/src/model"
	"github.com/sportivaid/go-template/util/httputil"
)

type Middleware struct {
	au AuthUsecase
}

func (m *Middleware) AuthUserToken(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		authHeader := c.GetHeader("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			log.Println(apperror.InvalidAuthToken)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, apperror.InvalidAuthToken)
			return
		}

		authUser, err := m.au.AuthenticateUserToken(bearerToken[1])
		if err != nil {
			log.Println(err)
			processTime := time.Now().Sub(startTime).Seconds()
			httputil.WriteErrorResponse(c, processTime, err)
			return
		}

		c.Set(model.AccountHeaderUserInfo, authUser)
		next(c)
	}
}

func NewMiddleware(au AuthUsecase) *Middleware {
	return &Middleware{
		au,
	}
}
